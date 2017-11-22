//SPDX-License-Identifier: Apache-2.0

var obj = require('../controllers/controller.js');

module.exports = function(app){

  app.get('/get_tuna/:id', function(req, res){
    obj.get_tuna(req, res);
  });
  app.get('/add_tuna/:tuna', function(req, res){
    obj.add_tuna(req, res);
  });
  app.get('/tuna_change_holder/:holder', function(req, res){
    obj.change_holder(req, res);
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
