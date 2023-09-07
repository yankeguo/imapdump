package imapdump

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/guoyk93/rg"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func DumpMessagePath(dir string, msg *imap.Message) string {
	year := msg.Envelope.Date.Format("2006")
	month := msg.Envelope.Date.Format("01")
	id := SanitizeMessageID(msg.Envelope.MessageId)
	if id == "" {
		// extract first from
		from := "unknown@unknown.com"
		if len(msg.Envelope.From) > 0 {
			from = strings.ToLower(msg.Envelope.From[0].Address())
		}
		// digest date, from addresses, to addesses, subject
		h := md5.New()
		_, _ = h.Write([]byte(strconv.FormatInt(msg.Envelope.Date.Unix(), 10)))
		for _, addr := range msg.Envelope.From {
			_, _ = h.Write([]byte(addr.Address()))
		}
		for _, addr := range msg.Envelope.To {
			_, _ = h.Write([]byte(addr.Address()))
		}
		_, _ = h.Write([]byte(msg.Envelope.Subject))
		digest := hex.EncodeToString(h.Sum(nil))
		// build a id
		id = msg.Envelope.Date.UTC().Format("20060102150405") + "-" + from + "-" + digest
	}
	return filepath.Join(dir, year, month, id+".eml")
}

func DumpMessagePathExisted(dir string, msg *imap.Message) (ok bool, err error) {
	file := DumpMessagePath(dir, msg)
	if _, err = os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	ok = true
	return
}

func DumpMessage(dir string, msg *imap.Message) (err error) {
	file := DumpMessagePath(dir, msg)

	if err = os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return
	}

	tmpFile := file + ".tmp"
	defer os.RemoveAll(tmpFile)

	r := msg.GetBody(&imap.BodySectionName{})

	if r == nil {
		log.Println("message does not have a body:", msg.Envelope.MessageId)
		return
	}

	var f *os.File
	if f, err = os.OpenFile(tmpFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0640); err != nil {
		return
	}

	_, err = io.Copy(f, r)

	_ = f.Close()

	if err != nil {
		return
	}

	if err = os.Rename(tmpFile, file); err != nil {
		return
	}

	return
}

type DumpAccountOptions struct {
	DisplayName string
	Dir         string
	Host        string
	Username    string
	Password    string
	Prefixes    []string
}

func DumpAccountMailbox(ctx context.Context, dir string, c *client.Client, mailbox string) (count int64, err error) {
	defer rg.Guard(&err)

	status := rg.Must(c.Select(mailbox, true))
	log.Printf("[%s]: %d total", status.Name, status.Messages)

	if status.Messages == 0 {
		return
	}

	seq := new(imap.SeqSet)
	seq.AddRange(1, status.Messages)

	chMsg := make(chan *imap.Message)
	chErr := make(chan error)

	go func() {
		chErr <- c.Fetch(seq, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, chMsg)
	}()

	var (
		seqDL = new(imap.SeqSet)
		numDL int64
	)

outerLoopMsg:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case err = <-chErr:
			if err != nil {
				return
			}
			break outerLoopMsg
		case msg := <-chMsg:
			if msg == nil {
				continue outerLoopMsg
			}
			ok := rg.Must(DumpMessagePathExisted(dir, msg))
			if ok {
				continue outerLoopMsg
			}
			numDL += 1
			seqDL.AddNum(msg.Uid)
		}
	}

	if seqDL.Empty() {
		return
	}

	log.Printf("[%s]: %d to download", mailbox, numDL)

	{
		chMsg := make(chan *imap.Message)
		chErr := make(chan error)

		var section imap.BodySectionName

		go func() {
			chErr <- c.UidFetch(seqDL, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, chMsg)
		}()

	outerLoopDump:
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case err = <-chErr:
				if err != nil {
					return
				}
				break outerLoopDump
			case msg := <-chMsg:
				if msg == nil {
					continue outerLoopDump
				}
				count += 1
				rg.Must0(DumpMessage(dir, msg))
				log.Printf("[%s]: %d/%d", mailbox, count, numDL)
			}
		}
	}
	return
}

func DumpAccount(ctx context.Context, opts DumpAccountOptions) (err error) {
	defer rg.Guard(&err)

	c := rg.Must(client.DialTLS(opts.Host, nil))
	defer c.Close()

	log.Println("dialed:", opts.Host)

	rg.Must0(c.Login(opts.Username, opts.Password))
	defer c.Logout()

	log.Println("signed in:", opts.Username)

	rg.Must0(ctx.Err())

	var mailboxes []string

	{
		chBoxes := make(chan *imap.MailboxInfo)
		chErr := make(chan error)

		go func() {
			chErr <- c.List("", "*", chBoxes)
		}()

	outerLoopBox:
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case err = <-chErr:
				if err != nil {
					return
				}
				break outerLoopBox
			case box := <-chBoxes:
				if box == nil {
					continue outerLoopBox
				}
				log.Printf("found [%s]", box.Name)
				for _, prefix := range opts.Prefixes {
					if strings.HasPrefix(box.Name, prefix) {
						log.Printf("matched [%s]", box.Name)
						mailboxes = append(mailboxes, box.Name)
						break
					}
				}
			}
		}
	}

	var count int64

	for _, mailbox := range mailboxes {
		rg.Must0(ctx.Err())

		count += rg.Must(DumpAccountMailbox(ctx, opts.Dir, c, mailbox))
	}

	log.Println("dumped:", count, "messages")

	return
}
