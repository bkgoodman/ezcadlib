package main

import (
	 "database/sql"
	"fmt"
	"strconv"
	"strings"
	"regexp"
	_ "github.com/mattn/go-sqlite3"

        "net/http"
        "net/http/cgi"
)

func replaceMap(s string,m *map[string]string) string {
	r := regexp.MustCompile("\\${[^}]*}")
	for x,i := range *m {
		s = strings.Replace(s,"${"+x+"}",i,-1)
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
<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ho+j7jyWK8fNQe+A12Hb8AhRq26LrZ/JpcUGGOn+Y7RsweNrtN/tE3MoK7ZeZDyx" crossorigin="anonymous"></script>`
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
	  <tbody> `
const TABLE_FOOTER = `
	  </tbody>
	</table>
	<a type="button" href="?new" class="btn btn-primary">New</a>
</div> `

const FORM = `
<form method="POST">
<input type="hidden" name="record" value="${recordno}" />
<button class="btn btn-primary" name="submit" value="update" type="submit">Update</button>
<button class="btn btn-primary" name="submit" value="savenew" type="submit">Save as New Copy</button>
   <div class="border rounded  container"> <!-- Master Container -->
	<div class="row g-2" class="form-control-sm" >
          <div class="col-md-4">
            <label for="validationServer01" class="form-label-sm form-label">Material</label>
            <input name="material" value="${material}" type="text" class="form-control-sm form-control" id="validationServer01" required="">
            <small class="form-control-sm" > e.g. "Aluminum"</small>
          </div>
          <div class="col-md-4">
            <label for="validationServer02" class="form-label-sm form-label">Operation</label>
            <input name="operation" value="${operation}" type="text" class="form-control-sm form-control" id="validationServer02" required="">
	    <small>e.g. "Mark, Black"</small>
          </div>
          <div class="col-md-4">
            <label for="validationServer03" class="form-label">User</label>
            <input name="user" value="${user}" type="text" class="form-control-sm form-control" id="validationServer03" required="">
	    <small>Your Name</small>
          </div>
          <div class="col-md-12">
            <label for="validationServer04" class="form-label">Comments</label>
            <input name="comments" value="${comments}" type="text" class="form-control-sm form-control" id="validationServer04">
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
            <input name="speed1" value="${speed1}" type="text" class="form-control-sm form-control" id="validationServer01" required="">
            <small> mm/Sec</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer02" class="form-label">Power</label>
            <input name="power1" value="${power1}" type="text" class="form-control-sm form-control" id="validationServer02" required="">
	    <small>Percent</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Frequency</label>
            <input name="frequency1" value="${frequency1}" type="text" class="form-control" id="validationServer03" required="">
	    <small>kHz</small>
          </div>

          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Pulse Width</label>
            <input name="pulsewidth1" value="${pulsewidth1}" type="text" class="form-control" id="validationServer03" >
	    <small>ns</small>
          </div>

	</div>
	<div class="row g-4">
		<label for="hatchType" class="form-label">Hatch Type</label>
		<div class="btn-group col-md-8" data-toggle="buttons">
			  <input name="hatch1" value="${hatchA1}" type="radio" class="btn btn-outline"><img src="static/hatch1.png" /></button>
			  <input name="hatch1" value="${hatchB1}" type="radio" class="btn btn-outline"><img src="static/hatch2.png" /></button>
			  <input name="hatch1" value="${hatchC1}" type="radio" class="btn btn-outline"><img src="static/hatch3.png" /></button>
			  <input name="hatch1" value="${hatchD1}" type="radio" class="btn btn-outline"><img src="static/hatch4.png" /></button>
		</div>
	</div>
	<div class="row">
          <div class="col-md-1">Angle</div>
          <div class="col-md-1">Loop Count</div>
          <div class="col-md-1">Line Space</div>
          <div class="col-md-1">Edge Offset</div>
          <div class="col-md-1">Start Offset</div>
          <div class="col-md-1">End Offset</div>
          <div class="col-md-1">Line Reduction</div>
          <div class="col-md-1">Loop Distance</div>
          <div class="col-md-2">Auto Rotate Hatch Angle</div>
  	</div>
	<div class="row">
          <div class="col-md-1">
            <input name="angle1" value="${angle1}" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small>Degrees</small>
	   </div>

          <div class="col-md-1">
            <input name="loopcount1" value="${loopcount1}" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Count</small>
          </div>

          <div class="col-md-1">
            <input name="linespace1" value="${linespace1}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>

          <div class="col-md-1">
            <input name="edgeoffset1" value="${edgeoffset1}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="startoffset1" value="${startoffset1}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="endoffset1" value="${endoffset1}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="linereduction1" value="${linereduction1}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="loopdistance1" value="${loopdistance1}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-2">
		  <input name="degrees1" value="${degrees1}" type="text" class="form-control" id="validationServer03" >
		  <input class="form-check-input" name="autorotate1" value="${autorotate1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  Enabled
          </div>



	 </div>


	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="markcountour1" value="${markcountour1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Mark Contour</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourA1" value="${contourA1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour1.png"></label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourB1" value="${contourB1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour2.png"></label>
		</div>
	  </div>
	</div>
	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="enable1" value="${enable1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="allcalc1" value="${allcalc1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">All Calc</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="followedgeonce1" value="${followedgeonce1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Follow Edge Once</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="crosshatch1" value="${crosshatch1}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Cross Hatch</label>
		</div>
	   </div>
         </div>
   </div>
</div></div> <!-- End Hatch Section -->




<p>
	<a class="btn btn-primary" name="secondCollapse" value="${secondCollapse}" data-toggle="collapse" href="#multiCollapseExample1" role="button" aria-expanded="false" aria-controls="multiCollapseExample1">Add Second Hatch</a>
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
            <input name="speed2" value="${speed2}" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small> mm/Sec</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer02" class="form-label">Power</label>
            <input name="power2" value="${power2}" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Percent</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Frequency</label>
            <input name="frequency2" value="${frequency2}" type="text" class="form-control" id="validationServer03" >
	    <small>kHz</small>
          </div>

          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Pulse Width</label>
            <input name="pulsewidth2" value="${pulsewidth2}" type="text" class="form-control" id="validationServer03" >
	    <small>ns</small>
          </div>

	</div>
	<div class="row g-4">
		<label for="hatchType" class="form-label">Hatch Type</label>
		<div class="btn-group col-md-8" data-toggle="buttons">
			  <input name="hatch2" value="${hatchA2}" type="radio" class="btn btn-outline"><img src="static/hatch1.png" /></button>
			  <input name="hatch2" value="${hatchB2}" type="radio" class="btn btn-outline"><img src="static/hatch2.png" /></button>
			  <input name="hatch2" value="${hatchC2}" type="radio" class="btn btn-outline"><img src="static/hatch3.png" /></button>
			  <input name="hatch2" value="${hatchD2}" type="radio" class="btn btn-outline"><img src="static/hatch4.png" /></button>
		</div>
	</div>
	<div class="row">
          <div class="col-md-1">Angle</div>
          <div class="col-md-1">Loop Count</div>
          <div class="col-md-1">Line Space</div>
          <div class="col-md-1">Edge Offset</div>
          <div class="col-md-1">Start Offset</div>
          <div class="col-md-1">End Offset</div>
          <div class="col-md-1">Line Reduction</div>
          <div class="col-md-1">Loop Distance</div>
          <div class="col-md-2">Auto Rotate Hatch Angle</div>
  	</div>
	<div class="row">
          <div class="col-md-1">
            <input name="angle2" value="${angle2}" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small>Degrees</small>
	   </div>

          <div class="col-md-1">
            <input name="loopcount2" value="${loopcount2}" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Count</small>
          </div>

          <div class="col-md-1">
            <input name="linespace2" value="${linespace2}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>

          <div class="col-md-1">
            <input name="edgeoffset2" value="${edgeoffset2}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="startoffset2" value="${startoffset2}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="endoffset2" value="${endoffset2}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="linereduction2" value="${linereduction2}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="loopdistance2" value="${loopdistance2}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-2">
            <input name="degrees2" value="${degrees2}" type="text" class="form-control" id="validationServer03" >
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="autorotate2" value="${autorotate2}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
	</div>



	 </div>


	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="markcountour2" value="${markcountour2}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Mark Contour</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourA2" value="${contourA2}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour1.png"></label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourB2" value="${contourB2}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour2.png"></label>
		</div>
	  </div>
	</div>
	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="enable2" value="${enable2}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="allcalc2" value="${allcalc2}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">All Calc</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="followedgeonce2" value="${followedgeonce2}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Follow Edge Once</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="crosshatch2" value="${crosshatch2}" type="checkbox" id="inlineCheckbox1" value="option1">
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
            <input name="speed3" value="${speed3}" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small> mm/Sec</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer02" class="form-label">Power</label>
            <input name="power3" value="${power3}" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Percent</small>
          </div>
          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Frequency</label>
            <input name="frequency3" value="${frequency3}" type="text" class="form-control" id="validationServer03" >
	    <small>kHz</small>
          </div>

          <div class="col-md-3">
            <label for="validationServer03" class="form-label">Pulse Width</label>
            <input name="pulsewidth3" value="${pulsewidth3}" type="text" class="form-control" id="validationServer03" >
	    <small>ns</small>
          </div>

	</div>
	<div class="row g-4">
		<label for="hatchType" class="form-label">Hatch Type</label>
		<div class="btn-group col-md-8" data-toggle="buttons">
			  <input name="hatch3" value="${hatch3}" value="${hatchA3}" type="radio" class="btn btn-outline"><img src="static/hatch1.png" /></button>
			  <input name="hatch3" value="${hatch3}" value="${hatchB3}" type="radio" class="btn btn-outline"><img src="static/hatch2.png" /></button>
			  <input name="hatch3" value="${hatch3}" value="${hatchC3}" type="radio" class="btn btn-outline"><img src="static/hatch3.png" /></button>
			  <input name="hatch3" value="${hatch3}" value="${hatchD3}" type="radio" class="btn btn-outline"><img src="static/hatch4.png" /></button>
		</div>
	</div>

	<div class="row">
          <div class="col-md-1">Angle</div>
          <div class="col-md-1">Loop Count</div>
          <div class="col-md-1">Line Space</div>
          <div class="col-md-1">Edge Offset</div>
          <div class="col-md-1">Start Offset</div>
          <div class="col-md-1">End Offset</div>
          <div class="col-md-1">Line Reduction</div>
          <div class="col-md-1">Loop Distance</div>
          <div class="col-md-2">Auto Rotate Hatch Angle</div>
  	</div>
	<div class="row">
          <div class="col-md-1">
            <input name="angle3" value="${angle3}" type="text" class="form-control-sm form-control" id="validationServer01" >
            <small>Degrees</small>
	   </div>

          <div class="col-md-1">
            <input name="loopcount3" value="${loopcount3}" type="text" class="form-control-sm form-control" id="validationServer02" >
	    <small>Count</small>
          </div>

          <div class="col-md-1">
            <input name="linespace3" value="${linespace3}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>

          <div class="col-md-1">
            <input name="edgeoffset3" value="${edgeoffset3}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="startoffset3" value="${startoffset3}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="endoffset3" value="${endoffset3}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="linereduction3" value="${linereduction3}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-1">
            <input name="loopdistance3" value="${loopdistance3}" type="text" class="form-control" id="validationServer03" >
	    <small>mm</small>
          </div>


          <div class="col-md-2">
            <input name="degrees3" value="${degrees3}" type="text" class="form-control" id="validationServer03" >
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="autorotate3" value="${autorotate3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
          </div>



	 </div>


	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="markcountour3" value="${markcountour3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Mark Contour</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourA3" value="${contourA3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour1.png"></label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="contourB3" value="${contourB3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1"><img src="static/contour2.png"></label>
		</div>
	  </div>
	</div>
	<div class="row">
          <div class="col-md-12">
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="enable3" value="${enable3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Enable</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="allcalc3" value="${allcalc3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">All Calc</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="followedgeonce3" value="${followedgeonce3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Follow Edge Once</label>
		</div>
		<div class="form-check form-check-inline">
		  <input class="form-check-input" name="crosshatch3" value="${crosshatch3}" type="checkbox" id="inlineCheckbox1" value="option1">
		  <label class="form-check-label" for="inlineCheckbox1">Cross Hatch</label>
		</div>
	   </div>
         </div>
   </div>
</div>
</div> <!-- Collase 2 -->
</div>  <!-- Collapse 1 -->



</div> <!-- Master Container -->
</form> `

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
