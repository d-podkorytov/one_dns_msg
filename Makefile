all: one_dns_msg

one_dns_msg: one_dns_msg.go Makefile
	go build one_dns_msg.go