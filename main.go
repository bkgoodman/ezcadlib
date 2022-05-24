package main

import (
	 "database/sql"
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3"

        "net/http"
        "net/http/cgi"
)
func checkErr(err error) {
	if err != nil {
	    panic(err)
	}
}

const HEADER = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Bootstrap demo</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-0evHe/X+R7YkIZDRvuzKMRqM+OrBnVFBL6DOitfPri4tjfHxaWutUpFmBp4vmVor" crossorigin="anonymous">
    <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.0/css/all.css" integrity="sha384-lZN37f5QGtY3VHgisS14W3ExzMWZxybE1SJSEsQp9S+oqd12jhcu+A56Ebc1zFSJ" crossorigin="anonymous">
  </head>
  <body>
   <main class="flex-shrink-0">
    <h1>EZCAD MOPA Settings</h1>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/js/bootstrap.bundle.min.js" integrity="sha384-pprn3073KE6tl6bjs2QrFaJGz5/SUsLqktiwsUTF55Jfv3qYSDhgCecCxMW52nD2" crossorigin="anonymous"></script>


    <script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ho+j7jyWK8fNQe+A12Hb8AhRq26LrZ/JpcUGGOn+Y7RsweNrtN/tE3MoK7ZeZDyx" crossorigin="anonymous"></script>

`
const FOOTER = `
</pre></main></body></html>`

const TABLE_HEADER = `
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
	  <tbody>
`
const TABLE_FOOTER = `
	  </tbody>
	</table>
	<a type="button" href="?new" class="btn btn-primary">New</a>
</div>
`

const FORM = `
<form method="POST">
<button class="btn btn-primary" type="submit">Update</button>
<button class="btn btn-primary" type="submit">Save as New Copy</button>
   <div class="border rounded  container"> <!-- Master Container -->
	<div class="row g-2" class="form-control-sm" >
          <div class="col-md-4">
            <label for="validationServer01" class="form-label-sm form-label">Material</label>
            <input name="material" type="text" class="form-control-sm form-control" id="validationServer01" required="">
            <small class="form-control-sm" > e.g. "Aluminum"</small>
          </div>
          <div class="col-md-4">
            <label for="validationServer02" class="form-label-sm form-label">Operation</label>
            <input name="operation" type="text" class="form-control-sm form-control" id="validationServer02" required="">
	    <small>e.g. "Mark, Black"</small>
          </div>
          <div class="col-md-4">
            <label for="validationServer03" class="form-label">User</label>
            <input name="user" type="text" class="form-control-sm form-control" id="validationServer03" required="">
	    <small>Your Name</small>
          </div>
          <div class="col-md-12">
            <label for="validationServer04" class="form-label">Comments</label>
            <input name="comments" type="text" class="form-control-sm form-control" id="validationServer04">
          </div>

	</div>





     <!-- hatch -->
<div class="row container">
  <div class="col">
    <div class="uncollapse multi-collapse" id="multiCollapseExample0">
      <dir class="card card-body">
	<div class="row g-4">
	<h3>Hatch</h3>
	</div>
	<div class="row g-4">
          <div class="col-md-3">
            <label for="validationServer01" class="form-label">Speed</label>
            <input name="speed1" type="text" class="form-control-sm form-control" id="validationServer01" required="">
            <small> mm/Sec</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer02" class="form-label">Power</label>
            <input name="power1" type="text" class="form-control-sm form-control" id="validationServer02" required="">
	    <small>Percent</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Frequency</label>
            <input name="frequency1" type="text" class="form-control" id="validationServer03" required="">
	    <small>kHz</small>
          </div>

          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Pulse Width</label>
            <input name="pulsewidth1" type="text" class="form-control" id="validationServer03" required="">
	    <small>ns</small>
          </div>

	</div>
	<div class="row g-4">
		<label for="hatchType" class="form-label">Hatch Type</label>
		<div class="btn-group col-md-8" role="group" aria-label="Hatch Type">
			  <button name="hatchA1" type="button" class="btn btn-outline"><img src="static/hatch1.png" /></button>
			  <button name="hatchB1" type="button" class="btn btn-outline"><img src="static/hatch2.png" /></button>
			  <button name="hatchC1" type="button" class="btn btn-outline"><img src="static/hatch3.png" /></button>
			  <button name="hatchD1" type="button" class="btn btn-outline"><img src="static/hatch4.png" /></button>
		</div>
	</div>
	<div class="row">
          <div class="col-md-1">
            <label for="validationServer01" class="form-label">Angle</label>
            <input name="angle1" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small>Degrees</small>
	   </div>

          <div class="col-md-1">
            <label for="validationServer02" class="form-label">Loop Count</label>
            <input name="loopcount1" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Count</small>
          </div>

          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Line Space</label>
            <input name="linespace1" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>

          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Edge Offset</label>
            <input name="edgeoffset1" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Start Offset</label>
            <input name="startoffset1" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">End Offset</label>
            <input name="endoffset1" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Line Reduction</label>
            <input name="linereduction1" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Loop Distance</label>
            <input name="loopdistance1" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="autorotate1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Auto Rotate Hatch Angle</label>
            <label for="validationServer03" class="form-label">Degrees</label>
            <input name="degrees1" type="text" class="form-control" id="validationServer03" >
		</div>



	 </div>


	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="markcountour1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Mark Contour</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourA1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour1.png"></label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourB1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour2.png"></label>
		</div>
	  </div>
	</div>
	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="enable1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="allcalc1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">All Calc</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="followedgeonce1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Follow Edge Once</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="crosshatch1" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Cross Hatch</label>
		</div>
	   </div>
         </div>
   </div>
</div></div> <!-- End Hatch Section -->

   <p>
  <a class="btn btn-primary" name="secondCollapse" data-toggle="collapse" href="#multiCollapseExample1" role="button" aria-expanded="false" aria-controls="multiCollapseExample1">Add Second Hatch</a>
</p>


     <!-- hatch -->
<div class="collapse multi-collapse" id="multiCollapseExample1">
  <div class="row container">
    <div class="col">
      <dir class="card card-body">
	<div class="row g-4">
	<h3>Hatch 2</h3>
	</div>
	<div class="row g-4">
          <div class="col-md-3">
            <label for="validationServer01" class="form-label">Speed</label>
            <input name="speed2" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small> mm/Sec</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer02" class="form-label">Power</label>
            <input name="power2" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Percent</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Frequency</label>
            <input name="frequency2" type="text" class="form-control" id="validationServer03" >
	    <small>kHz</small>
          </div>

          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Pulse Width</label>
            <input name="pulsewidth2" type="text" class="form-control" id="validationServer03" >
	    <small>ns</small>
          </div>

	</div>
	<div class="row g-4">
		<label for="hatchType" class="form-label">Hatch Type</label>
		<div class="btn-group col-md-8" role="group" aria-label="Hatch Type">
			  <button name="hatchA2" type="button" class="btn btn-outline"><img src="static/hatch1.png" /></button>
			  <button name="hatchB2" type="button" class="btn btn-outline"><img src="static/hatch2.png" /></button>
			  <button name="hatchC2" type="button" class="btn btn-outline"><img src="static/hatch3.png" /></button>
			  <button name="hatchD2" type="button" class="btn btn-outline"><img src="static/hatch4.png" /></button>
		</div>
	</div>
	<div class="row">
          <div class="col-md-1">
            <label for="validationServer01" class="form-label">Angle</label>
            <input name="angle2" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small>Degrees</small>
	   </div>

          <div class="col-md-1">
            <label for="validationServer02" class="form-label">Loop Count</label>
            <input name="loopcount2" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Count</small>
          </div>

          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Line Space</label>
            <input name="linespace2" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>

          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Edge Offset</label>
            <input name="edgeoffset2" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Start Offset</label>
            <input name="startoffset2" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">End Offset</label>
            <input name="endoffset2" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Line Reduction</label>
            <input name="linereduction2" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Loop Distance</label>
            <input name="loopdistance2" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="autorotate2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Auto Rotate Hatch Angle</label>
            <label for="validationServer03" class="form-label">Degrees</label>
            <input name="degrees2" type="text" class="form-control" id="validationServer03" >
		</div>



	 </div>


	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="markcountour2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Mark Contour</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourA2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour1.png"></label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourB2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour2.png"></label>
		</div>
	  </div>
	</div>
	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="enable2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="allcalc2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">All Calc</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="followedgeonce2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Follow Edge Once</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="crosshatch2" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Cross Hatch</label>
		</div>
	   </div>

   <p>
   </div> 
       </div> <!-- Column -->
</div> <!-- Row -->
  <a class="btn btn-primary" name=thirdCollapse" data-toggle="collapse" href="#multiCollapseExample2" role="button" aria-expanded="false" aria-controls="multiCollapseExample2">Add Third Hatch</a>
</p>





     <!-- hatch -->
<div class="collapse multi-collapse" id="multiCollapseExample2">
  <div class="row container">
    <div class="col">
      <dir class="card card-body">
	<div class="row g-4">
	<h3>Hatch 3</h3>
	</div>
	<div class="row g-4">
          <div class="col-md-3">
            <label for="validationServer01" class="form-label">Speed</label>
            <input name="speed3" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small> mm/Sec</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer02" class="form-label">Power</label>
            <input name="power3" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Percent</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Frequency</label>
            <input name="frequency3" type="text" class="form-control" id="validationServer03" >
	    <small>kHz</small>
          </div>

          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Pulse Width</label>
            <input name="pulsewidth3" type="text" class="form-control" id="validationServer03" >
	    <small>ns</small>
          </div>

	</div>
	<div class="row g-4">
		<label for="hatchType" class="form-label">Hatch Type</label>
		<div class="btn-group col-md-8" role="group" aria-label="Hatch Type">
			  <button name="hatchA3" type="button" class="btn btn-outline"><img src="static/hatch1.png" /></button>
			  <button name="hatchB3" type="button" class="btn btn-outline"><img src="static/hatch2.png" /></button>
			  <button name="hatchC3" type="button" class="btn btn-outline"><img src="static/hatch3.png" /></button>
			  <button name="hatchD3" type="button" class="btn btn-outline"><img src="static/hatch4.png" /></button>
		</div>
	</div>
	<div class="row">
          <div class="col-md-1">
            <label for="validationServer01" class="form-label">Angle</label>
            <input name="angle3" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small>Degrees</small>
	   </div>

          <div class="col-md-1">
            <label for="validationServer02" class="form-label">Loop Count</label>
            <input name="loopcount3" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Count</small>
          </div>

          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Line Space</label>
            <input name="linespace3" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>

          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Edge Offset</label>
            <input name="edgeoffset3" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Start Offset</label>
            <input name="startoffset3" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">End Offset</label>
            <input name="endoffset3" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Line Reduction</label>
            <input name="linereduction3" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <label for="validationServer03" class="form-label">Loop Distance</label>
            <input name="loopdistance3" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="autorotate3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Auto Rotate Hatch Angle</label>
            <label for="validationServer03" class="form-label">Degrees</label>
            <input name="degrees3" type="text" class="form-control" id="validationServer03" >
		</div>



	 </div>


	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="markcountour3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Mark Contour</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourA3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour1.png"></label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourB3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour2.png"></label>
		</div>
	  </div>
	</div>
	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="enable3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="allcalc3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">All Calc</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="followedgeonce3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Follow Edge Once</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="crosshatch3" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Cross Hatch</label>
		</div>
	   </div>
         </div>
   </div>
</div>
</div> <!-- Collase 2 -->
</div>  <!-- Collapse 1 -->



</div> <!-- Master Container -->
</form>
`

func showForm(w http.ResponseWriter, r *http.Request, db *sql.DB,item int) {
	fmt.Fprintln(w,FORM)
}

func showTable(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fmt.Fprintln(w,TABLE_HEADER)
	row, err := db.Query("SELECT id,material,op,user FROM params;")
	checkErr(err)
	defer row.Close()
	for row.Next() {
		var id int
		var material string
		var op string
		var user string
		row.Scan(&id, &material, &op, &user)
		fmt.Fprintf(w,`
	    <tr>
	      <td><a type="button" href="?detail=%d" class="btn btn-info fas fa-eye"></a></td>
	      <td>%s</td>
	      <td>%s</td>
	      <td>%s</td>
	      <td></td>
	    </tr>`,id,material,op,user)
	}
	fmt.Fprintln(w,TABLE_FOOTER)
} 
func main (){
	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		defer func(){fmt.Fprintln(w, FOOTER)}()

		db, err := sql.Open("sqlite3", "./params.db")
		fmt.Fprintln(w,err)
		defer func() { db.Close() }()
		checkErr(err)

		fmt.Fprintln(w, HEADER)

		r.ParseForm()
		form := r.Form
		if (form.Get("detail") == "") {
			showTable(w,r,db)
		} else {
			item,_ := strconv.Atoi(form.Get("detail"))
			showForm(w,r,db,item)
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


	

	// insert
        stmt, err := db.Prepare("INSERT INTO params(material, op, user,comments) values(?,?,?,?)")
        checkErr(err)

        res, err := stmt.Exec("Aluminum","Engrave Black","Brad", "This is a test")
      	if (err!= nil) {
		fmt.Fprintln(w,"Insert error ",err)
	} else {

		id, err := res.LastInsertId()
		fmt.Fprintln(w,"Inserted id",id)
		checkErr(err)
	}

	})); err != nil {
		fmt.Println(err)
	}
}
