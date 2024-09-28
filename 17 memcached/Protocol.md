## Protocol

Clients of memcached communicate with server through TCP connections.
(A UDP interface is also available; details are below under "UDP
protocol.") A given running memcached server listens on some
(configurable) port; clients connect to that port, send commands to
the server, read responses, and eventually close the connection.

There is no need to send any command to end the session. A client may
just close the connection at any moment it no longer needs it. Note,
however, that clients are encouraged to cache their connections rather
than reopen them every time they need to store or retrieve data. This
is because memcached is especially designed to work very efficiently
with a very large number (many hundreds, more than a thousand if
necessary) of open connections. Caching connections will eliminate the
overhead associated with establishing a TCP connection (the overhead
of preparing for a new connection on the server side is insignificant
compared to this).

There are two kinds of data sent in the memcache protocol: text lines
and unstructured data. Text lines are used for commands from clients
and responses from servers. Unstructured data is sent when a client
wants to store or retrieve data. The server will transmit back
unstructured data in exactly the same way it received it, as a byte
stream. The server doesn't care about byte order issues in
unstructured data and isn't aware of them. There are no limitations on
characters that may appear in unstructured data; however, the reader
of such data (either a client or a server) will always know, from a
preceding text line, the exact length of the data block being
transmitted.

Text lines are always terminated by \r\n. Unstructured data is _also_
terminated by \r\n, even though \r, \n or any other 8-bit characters
may also appear inside the data. Therefore, when a client retrieves
data from a server, it must use the length of the data block (which it
will be provided with) to determine where the data block ends, and not
the fact that \r\n follows the end of the data block, even though it
does.

## Keys

Data stored by memcached is identified with the help of a key. A key
is a text string which should uniquely identify the data for clients
that are interested in storing and retrieving it. Currently the
length limit of a key is set at 250 characters (of course, normally
clients wouldn't need to use such long keys); the key must not include
control characters or whitespace.

## Commands

There are three types of commands.

Storage commands (there are six: "set", "add", "replace", "append"
"prepend" and "cas") ask the server to store some data identified by a
key. The client sends a command line, and then a data block; after
that the client expects one line of response, which will indicate
success or failure.

Retrieval commands ("get", "gets", "gat", and "gats") ask the server to
retrieve data corresponding to a set of keys (one or more keys in one
request). The client sends a command line, which includes all the
requested keys; after that for each item the server finds it sends to
the client one response line with information about the item, and one
data block with the item's data; this continues until the server
finished with the "END" response line.

All other commands don't involve unstructured data. In all of them,
the client sends one command line, and expects (depending on the
command) either one line of response, or several lines of response
ending with "END" on the last line.

A command line always starts with the name of the command, followed by
parameters (if any) delimited by whitespace. Command names are
lower-case and are case-sensitive.

## Meta Commands

Meta commands are a set of new ASCII-based commands. They follow the same
structure of the original commands but have a new flexible feature set.
Meta commands incorporate most features available in the binary protocol, as
well as a flag system to make the commands flexible rather than having a
large number of high level commands. These commands completely replace the
usage of basic Storage and Retrieval commands.

Meta flags are not to be confused with client bitflags, which is an opaque
number passed by the client. Meta flags change how the command operates, but
they are not stored in cache.

These work mixed with normal protocol commands on the same connection. All
existing commands continue to work. The meta commands will not replace
specialized commands that do not sit in the Storage or Retrieval categories
noted in the 'Commands' section above.

All meta commands follow a basic syntax:

<cm> <key> <datalen\*> <flag1> <flag2> <...>\r\n

Where <cm> is a 2 character command code.
<datalen> is only for commands with payloads, like set.

Responses look like:

<RC> <datalen\*> <flag1> <flag2> <...>\r\n

Where <RC> is a 2 character return code. The number of flags returned are
based off of the flags supplied.

<datalen> is only for responses with payloads, with the return code 'VA'.

Flags are single character codes, ie 'q' or 'k' or 'I', which adjust the
behavior of the command. If a flag requests a response flag (ie 't' for TTL
remaining), it is returned in the same order as they were in the original
command, though this is not strict.

Flags are single character codes, ie 'q' or 'k' or 'O', which adjust the
behavior of a command. Flags may contain token arguments, which come after the
flag and before the next space or newline, ie 'Oopaque' or 'Kuserkey'. Flags
can return new data or reflect information, in the same order they were
supplied in the request. Sending an 't' flag with a get for an item with 20
seconds of TTL remaining, would return 't20' in the response.

All commands accept a tokens 'P' and 'L' which are completely ignored. The
arguments to 'P' and 'L' can be used as hints or path specifications to a
proxy or router inbetween a client and a memcached daemon. For example, a
client may prepend a "path" in the key itself: "mg /path/foo v" or in a proxy
token: "mg foo Lpath/ v" - the proxy may then optionally remove or forward the
token to a memcached daemon, which will ignore them.

Syntax errors are handled the same as noted under 'Error strings' section
below.

For usage examples beyond basic syntax, please see the wiki:
https://github.com/memcached/memcached/wiki/MetaCommands

## Expiration times

Some commands involve a client sending some kind of expiration time
(relative to an item or to an operation requested by the client) to
the server. In all such cases, the actual value sent may either be
Unix time (number of seconds since January 1, 1970, as a 32-bit
value), or a number of seconds starting from current time. In the
latter case, this number of seconds may not exceed 60*60*24\*30 (number
of seconds in 30 days); if the number sent by a client is larger than
that, the server will consider it to be real Unix time value rather
than an offset from current time.

Note that a TTL of 1 will sometimes immediately expire. Time is internally
updated on second boundaries, which makes expiration time roughly +/- 1s.
This more proportionally affects very low TTL's.

## Error strings

Each command sent by a client may be answered with an error string
from the server. These error strings come in three types:

- "ERROR\r\n"

  means the client sent a nonexistent command name.

- "CLIENT_ERROR <error>\r\n"

  means some sort of client error in the input line, i.e. the input
  doesn't conform to the protocol in some way. <error> is a
  human-readable error string.

- "SERVER_ERROR <error>\r\n"

  means some sort of server error prevents the server from carrying
  out the command. <error> is a human-readable error string. In cases
  of severe server errors, which make it impossible to continue
  serving the client (this shouldn't normally happen), the server will
  close the connection after sending the error line. This is the only
  case in which the server closes a connection to a client.

In the descriptions of individual commands below, these error lines
are not again specifically mentioned, but clients must allow for their
possibility.

## Storage commands

First, the client sends a command line which looks like this:

<command name> <key> <flags> <exptime> <bytes> [noreply]\r\n
cas <key> <flags> <exptime> <bytes> <cas unique> [noreply]\r\n

- <command name> is "set", "add", "replace", "append" or "prepend"

  "set" means "store this data".

  "add" means "store this data, but only if the server _doesn't_ already
  hold data for this key".

  "replace" means "store this data, but only if the server _does_
  already hold data for this key".

  "append" means "add this data to an existing key after existing data".

  "prepend" means "add this data to an existing key before existing data".

  The append and prepend commands do not accept flags or exptime.
  They update existing data portions, and ignore new flag and exptime
  settings.

  "cas" is a check and set operation which means "store this data but
  only if no one else has updated since I last fetched it."

- <key> is the key under which the client asks to store the data

- <flags> is an arbitrary 16-bit unsigned integer (written out in
  decimal) that the server stores along with the data and sends back
  when the item is retrieved. Clients may use this as a bit field to
  store data-specific information; this field is opaque to the server.
  Note that in memcached 1.2.1 and higher, flags may be 32-bits, instead
  of 16, but you might want to restrict yourself to 16 bits for
  compatibility with older versions.

- <exptime> is expiration time. If it's 0, the item never expires
  (although it may be deleted from the cache to make place for other
  items). If it's non-zero (either Unix time or offset in seconds from
  current time), it is guaranteed that clients will not be able to
  retrieve this item after the expiration time arrives (measured by
  server time). If a negative value is given the item is immediately
  expired.

- <bytes> is the number of bytes in the data block to follow, _not_
  including the delimiting \r\n. <bytes> may be zero (in which case
  it's followed by an empty data block).

- <cas unique> is a unique 64-bit value of an existing entry.
  Clients should use the value returned from the "gets" command
  when issuing "cas" updates.

- "noreply" optional parameter instructs the server to not send the
  reply. NOTE: if the request line is malformed, the server can't
  parse "noreply" option reliably. In this case it may send the error
  to the client, and not reading it on the client side will break
  things. Client should construct only valid requests.

After this line, the client sends the data block:

<data block>\r\n

- <data block> is a chunk of arbitrary 8-bit data of length <bytes>
  from the previous line.

After sending the command line and the data block the client awaits
the reply, which may be:

- "STORED\r\n", to indicate success.

- "NOT_STORED\r\n" to indicate the data was not stored, but not
  because of an error. This normally means that the
  condition for an "add" or a "replace" command wasn't met.

- "EXISTS\r\n" to indicate that the item you are trying to store with
  a "cas" command has been modified since you last fetched it.

- "NOT_FOUND\r\n" to indicate that the item you are trying to store
  with a "cas" command did not exist.

## Retrieval command:

The retrieval commands "get" and "gets" operate like this:

get <key>_\r\n
gets <key>_\r\n

- <key>\* means one or more key strings separated by whitespace.

After this command, the client expects zero or more items, each of
which is received as a text line followed by a data block. After all
the items have been transmitted, the server sends the string

"END\r\n"

to indicate the end of response.

Each item sent by the server looks like this:

VALUE <key> <flags> <bytes> [<cas unique>]\r\n
<data block>\r\n

- <key> is the key for the item being sent

- <flags> is the flags value set by the storage command

- <bytes> is the length of the data block to follow, _not_ including
  its delimiting \r\n

- <cas unique> is a unique 64-bit integer that uniquely identifies
  this specific item.

- <data block> is the data for this item.

If some of the keys appearing in a retrieval request are not sent back
by the server in the item list this means that the server does not
hold items with such keys (because they were never stored, or stored
but deleted to make space for more items, or expired, or explicitly
deleted by a client).
