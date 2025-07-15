package sattva

import "github.com/hootuu/helix/components/sattva/channel"

type Identification = string

const IdNil = ""

const (
	Mobile   channel.Type = 100
	Password channel.Type = 101
	Email    channel.Type = 103
	DeviceID channel.Type = 104
	WeChat   channel.Type = 201
	AliPay   channel.Type = 202
)

type Channel = channel.Channel
