package model

type Postcard struct {
	Timestamp int64
	OrderNo   string
	Email     string `json:",omitempty"`
	Phone     uint64 `json:",omitempty"`
	IsSent    bool
}

func (d *Postcard) SetTimestamp(a int64) {
	d.Timestamp = a
}
func (d *Postcard) SetOrderNo(a string) {
	d.OrderNo = a
}
func (d *Postcard) SetEmail(a string) {
	d.Email = a
}
func (d *Postcard) SetPhone(a uint64) {
	d.Phone = a
}
