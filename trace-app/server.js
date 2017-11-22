//SPDX-License-Identifier: Apache-2.0

// nodejs server setup 

// call the packages we need
var express       = require('express');        // call express
var app           = express();                 // define our app using express
var bodyParser    = require('body-parser');
var http          = require('http')
var fs            = require('fs');
var Fabric_Client = require('fabric-client');
var path          = require('path');
var util          = require('util');
var os            = require('os');

// Load all of our middleware
// configure app to use bodyParser()
// this will let us get the data from a POST
// app.use(express.static(__dirname + '/client'));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

// instantiate the app
var app = express();

// this line requires and runs the code from our routes.js file and passes it app
<<<<<<< HEAD
require('./routes/routes.js')(app);
=======
require('./routes.js')(app);
>>>>>>> c7064045e57a27b1e36890c820fa0b1c46bb64b3

// set up a static file server that points to the "client" directory
app.use(express.static(path.join(__dirname, './client')));

// Save our port
<<<<<<< HEAD
var port = process.env.PORT || 3000;
=======
var port = process.env.PORT || 8000;
>>>>>>> c7064045e57a27b1e36890c820fa0b1c46bb64b3

// Start the server and listen on port 
app.listen(port, '0.0.0.0', function(){
  console.log("Live on port: " + port);
});
<<<<<<< HEAD
=======

>>>>>>> c7064045e57a27b1e36890c820fa0b1c46bb64b3
