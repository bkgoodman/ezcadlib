package main
type params struct {
	Id int64 `db:"id"`
	Material string  `db:"material"`
	Op string `db:"op"`
	User string `db:"user"`
	Comments string `db:"comments"`
}

type  hatches struct {
	id int64 `db:"id"`
	sequence int `db:"sequence"`
	param int `db:"param"`
	edgeoffset int `db:"edgeoffset"`
	loopcount int `db:"loopcount"`
	startoffset int `db:"startoffset"`
	angle int `db:"angle"`
	loopdistance int `db:"loopdistance"`
	frequency int `db:"frequency"`
	linespace int `db:"linespace"`
	speed int `db:"speed"`
	endoffset int `db:"endoffset"`
	linereduction int `db:"linereduction"`
	power int `db:"power"`
	pulsewidth int `db:"pulsewidth"`
	degrees int `db:"degrees"`
	hatch int `db:"hatch"`

	markcontour bool `db:"markcontour"`
	contoura bool `db:"contoura"`
	contourb bool `db:"contourb"`
	crosshatch bool `db:"crosshatch"`
	enable bool `db:"enable"`
	allcalc bool `db:"allcalc"`
	followedgeonce bool `db:"followedgeonce"`
	autorotate bool `db:"autorotate"`
}
