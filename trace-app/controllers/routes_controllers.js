//SPDX-License-Identifier: Apache-2.0

var obj = require('../controllers/controller.js');

module.exports = function(app){

  app.get('/register', function(req, res) {
    console.log('register')
    //res.send('register');
    res.sendFile('./views/register.html')
  });
  app.get('/query_with_id', function(req,res){
    res.send('query_with_id');
  });
  app.get('/update_butcher', function(req,res){
    res.send('update_butcher');
  });
  app.get('/update_package', function(req,res){
    res.send('update_package');
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
