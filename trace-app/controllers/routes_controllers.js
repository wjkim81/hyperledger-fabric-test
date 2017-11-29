//SPDX-License-Identifier: Apache-2.0

var obj = require('../controllers/controller.js');
var path          = require('path');

module.exports = function(app){

  app.get('/register', function(req, res) {
    res.sendFile(path.join(__dirname, '../views', 'register.html'));
  });
  app.get('/query', function(req,res){
    res.sendFile(path.join(__dirname, '../views', 'query.html'));
  });
  app.get('/update_butchery', function(req,res){
    res.sendFile(path.join(__dirname, '../views', 'butchery.html'));
  });
  app.get('/update_process', function(req,res){
    res.sendFile(path.join(__dirname, '../views', 'package.html'));
  });


  app.get('/get_cattle/:id', function(req, res){
    obj.get_cattle(req, res);
  });
  app.get('/add_cattle/:cattle', function(req, res){
    obj.add_cattle(req, res);
  });
  app.get('/get_all_cattle', function(req, res){
    obj.get_all_cattle(req, res);
  });
  app.get('/register_cattle/:cattle', function(req, res){
    obj.register_cattle(req, res);
  });
  app.get('/update_butchery_info/:butchery_info', function(req, res){
    obj.update_butchery_info(req, res);
  });
  app.get('/update_process_info/:process_info', function(req, res){
    obj.update_process_info(req, res);
  });
}
