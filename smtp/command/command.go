package command

type Name string

const (
	HELO Name = "HELO"
	EHLO Name = "EHLO"
	MAIL Name = "MAIL FROM"
	RCPT Name = "RCPT TO"
	DATA Name = "DATA"
	QUIT Name = "QUIT"
	RSET Name = "RSET"
	NOOP Name = "NOOP"
	VRFY Name = "VRFY"
	EXPN Name = "EXPN"
)

type Command struct {
	Name Name
	Args []string
}
