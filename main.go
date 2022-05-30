package main

import (
	//"database/sql"
	"fmt"
	"strconv"
	"strings"
	"log"
	"os"
	"regexp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"

	"net/url"
	"html"
	"reflect"
        "net/http"
        "net/http/cgi"
)

var checkboxes = []string{
	"markcontour",
	"contoura",
	"contourb",
	"crosshatch",
	"enable",
	"allcalc",
	"followedgeonce",
	"autorotate"}

// Take fields in a struct and stick them in a map
func structToMap(x any,m map[string]string,suffix string) {
	u:= reflect.ValueOf(x)
	t:= u.Type() // t := reflect.TypeOf(u)
	v := u //v := reflect.ValueOf(&u).Elem()

	//fmt.Printf("\nV is %T %+v\n", v, v)


	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		//fmt.Printf("\nField Name %s %s TAG:%s %T %q\n", field.Name, field.Type, field.Tag.Get("mytag"), value, value)
		if (value.Kind() == reflect.Int64)  {
			m[field.Name+suffix] = fmt.Sprintf("%d",value.Int())
		} else if (value.Kind() == reflect.Int)  {
			m[field.Name+suffix] = fmt.Sprintf("%d",value.Int())

		} else if (value.Kind() == reflect.Bool)  {
			if (value.Bool()) {
				m[field.Name+suffix] = "true"
			} else {
				m[field.Name+suffix] = "false"
			}

		} else {
			m[field.Name+suffix] = value.String()
		}
	}
}

// Take variables in the map and replace them with ${field} in string
func replaceMap(s string,m *map[string]string) string {
	r := regexp.MustCompile("\\${[^}]*}")
	for x,i := range *m {
		s = strings.Replace(s,"${"+html.EscapeString(x)+"}",i,-1)
	}
	// Remove missing parameters
	s = r.ReplaceAllString(s,"")
	return s
}

func checkErr(err error) {
	if err != nil {
	    panic(err)
	}
}

const FOOTER = `
</pre></main></body></html>`

const TABLE_HEADER = `
<form method="POST">
<div class="container">
	<table class="table table-striped">
	  <thead>
	    <tr>
	      <th scope="col"></th>
	      <th scope="col">Material</th>
	      <th scope="col">Operation</th>
	      <th scope="col">User</th>
	    </tr>
	  </thead>
	  <tbody> `
const TABLE_FOOTER = `
	  </tbody>
	</table>
</div></form> `


func getFile(fn string) string {
	f,_:=os.Open(fn)
	defer func(){f.Close()}()

	s,_ := f.Stat()
	l := s.Size()
	b := make([]byte,l)
	f.Read(b)
	return string(b)
}

func showForm(w http.ResponseWriter, r *http.Request, db *sqlx.DB,item int) {
	fmt.Fprintln(w,"<!--")
	stmt, err := db.Preparex("SELECT * FROM params WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	m := make(map[string]string)
	var p params
	err = stmt.QueryRowx(item).StructScan(&p)
	if err != nil {
		fmt.Fprintln(w,err)
	}

	structToMap(p,m,"")
	m["collapse2"]="collapse"
	m["collapse3"]="collapse"
	rows,err2 :=db.Query(fmt.Sprintf(`SELECT 
		id,sequence, param, edgeoffset, loopcount, startoffset, angle, loopdistance, frequency,
		linespace, speed, endoffset, linereduction, power,
		pulsewidth, degrees, hatch, markcontour, contoura,
		contourb, crosshatch, enable, allcalc, followedgeonce, autorotate
	FROM hatches WHERE id = %d ORDER BY sequence;`,item))
	if err2 != nil {
		fmt.Fprintln(w,err2)
	} else {
		fmt.Fprintf(w,"%d HATCH ROWS %+v\n",item,rows)
		for rows.Next() {
			//rows.Scan
			fmt.Fprintf(w,"HATCH ROW %+v\n",rows)
			var h hatches
			_ = rows.Scan(
				&h.id,
				&h.sequence,
				&h.param,
				&h.edgeoffset,
				&h.loopcount,
				&h.startoffset,
				&h.angle,
				&h.loopdistance,
				&h.frequency,
				&h.linespace,
				&h.speed,
				&h.endoffset,
				&h.linereduction,
				&h.power,
				&h.pulsewidth,
				&h.degrees,
				&h.hatch,

				&h.markcontour,
				&h.contoura,
				&h.contourb,
				&h.crosshatch,
				&h.enable,
				&h.allcalc,
				&h.followedgeonce,
				&h.autorotate)
			fmt.Fprintf(w,"Hatch %+v\n",h)
			m[fmt.Sprintf("collapse%d",h.sequence)]="uncollapse"
			structToMap(h,m,strconv.Itoa(h.sequence))
			// Intepret hatch selector/radio
			hname := fmt.Sprintf("hatch%c%d",'@'+h.hatch,h.sequence)
			m[hname]="checked"
			fmt.Println("HATCH SET",hname,"TO",m[hname])
			for _,x := range checkboxes  {
				hname = fmt.Sprintf("%s%d",x,h.sequence)
				if (m[hname] == "true") {
					m[hname]="checked"
				} else {
					m[hname]=""
				}
				fmt.Println("SET",hname,"TO",m[hname])
			}
		}
	}
	fmt.Fprintf(w,"Map is  %+v\n",m)
	fmt.Fprintln(w,"-->")
	form := getFile("form.html")
	fmt.Fprintln(w,replaceMap(form,&m))
	fmt.Fprintf(w,"%+v",m)

}

// Display a "new" form form
func postNew(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	m := make(map[string]string)
	m["Id"]="New"
	m["updateButtonModifier"]="style='display:none'"
	form := getFile("form.html")
	fmt.Fprintln(w,replaceMap(form,&m))

}

// replace an existing param
func saveParam(oid int64,form url.Values, db *sqlx.DB) {
	// insert
        stmt, err := db.Prepare("INSERT INTO params(material, op, user,comments) values(?,?,?,?)")
        checkErr(err)

	var id int64
        res, err := stmt.Exec(
		form.Get("material"),
		form.Get("operation"),
		form.Get("user"),
		form.Get("comments"))
      	if (err!= nil) {
		fmt.Println("Insert error ",err)
	} else {

		id, _ = res.LastInsertId()
		fmt.Println("<p>Inserted id",id,"</p>")
		checkErr(err)
	}
}

func saveHatch(oid int64,hatchno int,form url.Values, db *sqlx.DB) {
	var h hatches


	h.id=oid
	h.sequence=hatchno
	h.param,_ = strconv.Atoi(form.Get("param"+strconv.Itoa(hatchno)))
	h.edgeoffset,_ = strconv.Atoi(form.Get("edgeoffset"+strconv.Itoa(hatchno)))
	h.loopcount,_ = strconv.Atoi(form.Get("loopcount"+strconv.Itoa(hatchno)))
	h.startoffset,_ = strconv.Atoi(form.Get("startoffset"+strconv.Itoa(hatchno)))
	h.angle,_ = strconv.Atoi(form.Get("angle"+strconv.Itoa(hatchno)))
	h.loopdistance,_ = strconv.Atoi(form.Get("loopdistance"+strconv.Itoa(hatchno)))
	h.frequency,_ = strconv.Atoi(form.Get("frequency"+strconv.Itoa(hatchno)))
	h.linespace,_ = strconv.Atoi(form.Get("linespace"+strconv.Itoa(hatchno)))
	h.speed,_ = strconv.Atoi(form.Get("speed"+strconv.Itoa(hatchno)))
	h.endoffset,_ = strconv.Atoi(form.Get("endoffset"+strconv.Itoa(hatchno)))
	h.linereduction,_ = strconv.Atoi(form.Get("linereduction"+strconv.Itoa(hatchno)))
	h.power,_ = strconv.Atoi(form.Get("power"+strconv.Itoa(hatchno)))
	h.pulsewidth,_ = strconv.Atoi(form.Get("pulsewidth"+strconv.Itoa(hatchno)))
	h.degrees,_ = strconv.Atoi(form.Get("degrees"+strconv.Itoa(hatchno)))
	h.hatch,_ = strconv.Atoi(form.Get("hatch"+strconv.Itoa(hatchno)))

	h.markcontour = form.Get("markcontour"+strconv.Itoa(hatchno))=="true"
	h.contoura = form.Get("contoura"+strconv.Itoa(hatchno))=="true"
	h.contourb = form.Get("contourb"+strconv.Itoa(hatchno))=="true"
	h.crosshatch = form.Get("crosshatch"+strconv.Itoa(hatchno))=="true"
	h.enable = form.Get("enable"+strconv.Itoa(hatchno))=="true"
	h.allcalc = form.Get("allcalc"+strconv.Itoa(hatchno))=="true"
	h.followedgeonce = form.Get("followedgeonce"+strconv.Itoa(hatchno))=="true"
	h.autorotate = form.Get("autorotate"+strconv.Itoa(hatchno))=="true"

	q := `INSERT INTO hatches ( 
	id,sequence, param, edgeoffset, loopcount, startoffset, angle, loopdistance, frequency,
	linespace, speed, endoffset, linereduction, power,
	pulsewidth, degrees, hatch, markcontour, contoura,
	contourb, crosshatch, enable, allcalc, followedgeonce, autorotate)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_,err := db.Exec(q,
		h.id,h.sequence,
		h.param,
		h.edgeoffset,
		h.loopcount,
		h.startoffset,
		h.angle,
		h.loopdistance,
		h.frequency,
		h.linespace,
		h.speed,
		h.endoffset,
		h.linereduction,
		h.power,
		h.pulsewidth,
		h.degrees,
		h.hatch,

		h.markcontour,
		h.contoura,
		h.contourb,
		h.crosshatch,
		h.enable,
		h.allcalc,
		h.followedgeonce,
		h.autorotate)
      	if (err!= nil) {
		fmt.Println("Hatch Insert error ",err)
	}
	fmt.Printf("INSERTING %+v\n",h)
}

// Update an EXISTING record
func updateRecord(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	form := r.Form
	id,err := strconv.ParseInt(form.Get("Id"),10,64)
	fmt.Fprintln(w, "<h2>Update",id,"</h2>")
	m := make(map[string]string)
	_ = err
	_ = m
        stmt, err := db.Prepare("UPDATE params SET material =?, op= ?, user=?,comments=? WHERE id = ?")
        checkErr(err)

        res, err := stmt.Exec(
		form.Get("Material"),
		form.Get("Op"),
		form.Get("User"),
		form.Get("Comments"),id)
	_ = res
	_,err = db.Exec("DELETE FROM hatches where id = ?",id)
	/* Save Hatches */
	saveHatch(id,1,form,db)
	if (form.Get("hatchForm2")=="true") {
		fmt.Fprintln(w,"<!-- Add Second Hatch -->")
		saveHatch(id,2,form,db)
	}
	if ((form.Get("hatchForm2")=="true") && (form.Get("hatchForm3") == "true")) {
		fmt.Fprintln(w,"<!-- Add Third Hatch -->")
		saveHatch(id,3,form,db)
	}
}

// Save a NEW parameter
func saveNew(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	fmt.Fprintln(w,"<p>Saving New</p>")

	// insert
        stmt, err := db.Prepare("INSERT INTO params(material, op, user,comments) values(?,?,?,?)")
        checkErr(err)

	form := r.Form
	var id int64
        res, err := stmt.Exec(
		form.Get("Material"),
		form.Get("Op"),
		form.Get("User"),
		form.Get("Comments"))
      	if (err!= nil) {
		fmt.Fprintln(w,"Insert error ",err)
	} else {

		id, _ = res.LastInsertId()
		fmt.Fprintln(w,"<!-- Inserted id",id,"-->")
		checkErr(err)
	}

	/* Save Hatches */
	saveHatch(id,1,form,db)
	if (form.Get("hatchForm2")=="true") {
		fmt.Fprintln(w,"<!-- Add Second Hatch -->")
		saveHatch(id,2,form,db)
	}
	if ((form.Get("hatchForm2")=="true") && (form.Get("hatchForm3") == "true")) {
		fmt.Fprintln(w,"<!-- Add Third Hatch -->")
		saveHatch(id,3,form,db)
	}
}

func showTable(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	fmt.Fprintln(w,TABLE_HEADER)
	row, err := db.Queryx("SELECT * FROM params;")
	checkErr(err)
	defer row.Close()
	for row.Next() {
		var p params
		row.StructScan(&p)
		fmt.Fprintf(w,`
	    <tr>
	      <td>
	      	<button type="submit" name="viewRecord" value="%d" class="btn btn-info fas fa-eye"></button>
		</td>
	      <td>%s</td>
	      <td>%s</td>
	      <td>%s</td>
	      <td></td>
	    </tr>`,p.Id,p.Material,p.Op,p.User)
	}
	fmt.Fprintln(w,TABLE_FOOTER)
} 
func main (){
	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		defer func(){fmt.Fprintln(w, FOOTER)}()

		db, err := sqlx.Open("sqlite3", "./params.db")
		fmt.Fprintln(w,err)
		defer func() { db.Close() }()
		checkErr(err)

		fmt.Fprintln(w, getFile("header.html"))

		r.ParseForm()
		form := r.Form
		switch {
			case form.Get("submit") == "update":
				updateRecord(w,r,db)
			case form.Get("submit") == "savenew":
				fmt.Fprintln(w, "<h2>Save/New</h2>")
				saveNew(w,r,db)
			case form.Get("newSubmit") == "new":
				fmt.Fprintln(w, "<h2>Create New Entry</h2>")
				postNew(w,r,db)
			case form.Get("viewRecord") != "":
				item,_ := strconv.Atoi(form.Get("viewRecord"))
				showForm(w,r,db,item)
			default:
				showTable(w,r,db)
		}


		fmt.Fprintln(w,"<pre>")
		fmt.Fprintln(w, "Method:", r.Method)
		fmt.Fprintln(w, "URL:", r.URL.String())
		query := r.URL.Query()
		for k := range query {
			fmt.Fprintln(w, "Query", k+":", query.Get(k))
		}
		for k := range form {
			fmt.Fprintln(w, "Form", k+":", form.Get(k))
		}
		post := r.PostForm
		for k := range post {
			fmt.Fprintln(w, "PostForm", k+":", post.Get(k))
		}
		fmt.Fprintln(w, "RemoteAddr:", r.RemoteAddr)
		if referer := r.Referer(); len(referer) > 0 {
			fmt.Fprintln(w, "Referer:", referer)
		}
		if ua := r.UserAgent(); len(ua) > 0 {
			fmt.Fprintln(w, "UserAgent:", ua)
		}
		for _, cookie := range r.Cookies() {
			fmt.Fprintln(w, "Cookie", cookie.Name+":", cookie.Value, cookie.Path, cookie.Domain, cookie.RawExpires)
		}


	}) ); err != nil {
		fmt.Println(err)
	}
}
