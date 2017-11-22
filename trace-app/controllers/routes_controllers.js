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
  app.get('/update_butcher', function(req,res){
    res.sendFile(path.join(__dirname, '../views', 'butcher.html'));
  });
  app.get('/update_package', function(req,res){
    res.sendFile(path.join(__dirname, '../views', 'package.html'));
  });


  app.get('/get_cow/:id', function(req, res){
    obj.get_cow(req, res);
  });
  app.get('/add_cow/:cow', function(req, res){
    obj.add_cow(req, res);
  });
  app.get('/get_all_cow', function(req, res){
    obj.get_all_cow(req, res);
  });
  app.get('/register_cow/:cow', function(req, res){
    obj.register_cow(req, res);
  });
  app.get('/update_slaughter_info_cow/:slaughter_info', function(req, res){
    obj.update_slaughter_info_cow(req, res);
  });
  app.get('/update_package_info_cow/:package_info', function(req, res){
    obj.update_package_info_cow(req, res);
  });
  app.get('/cow_change_holder/:holder', function(req, res){
    obj.change_holder(req, res);
  });
}
