// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

    $("#success_holder").hide();
    $("#success_create").hide();
    $("#error_holder").hide();
    $("#error_query").hide();

    // Create angular function for traceability system
    $scope.queryAllCow = function(){

        appFactory.queryAllCow(function(data){
            var array = [];
            for (var i = 0; i < data.length; i++){
                //parseInt(data[i].Key);
                //data[i].Record.Key = parseInt(data[i].Key);
                data[i].Record.Key = data[i].Key;
                array.push(data[i].Record);
            }
            array.sort(function(a, b) {
                return parseFloat(a.Key) - parseFloat(b.Key);
            });
            $scope.all_cow = array;
        });
    }

    $scope.queryCow = function(){

        var id = $scope.trace_id;

        appFactory.queryCow(id, function(data){
            $scope.query_cow = data;

            if ($scope.query_cow == "Could not locate cow"){
                $("#error_query").show();
            } else{
                $("#error_query").hide();
            }

            //var array = [];
            //for (var i=0; i < data.package_info.length; i++){
            //    array.push(data.package_info[i].Record);                
            //}
            //array.sort(function(a, b) {
            //    return parseFloat(a.Key) - parseFloat(b.Key);
            //}
            //$scope.all_package_info = array;
            $scope.all_package_info = data.package_info; 
        });
    }

    $scope.registerCow = function(){

        appFactory.registerCow($scope.cow, function(data){
            $scope.create_cow = data;
            $("#success_create").show();
        });
    }

    $scope.updateSlaughterInfoCow = function(){

        appFactory.updateSlaughterInfoCow($scope.slaughter_info, function(data){
            $scope.create_cow = data;
            $("#success_create").show();
        });
    }

    $scope.updatePackageInfoCow = function(){

        appFactory.updatePackageInfoCow($scope.package_info, function(data){
            $scope.create_cow = data;
            $("#success_create").show();
        });
    }

    $scope.recordCow = function(){

        appFactory.recordCow($scope.cow, function(data){
            $scope.create_cow = data;
            $("#success_create").show();
        });
    }   

    $scope.changeHolder = function(){

        appFactory.changeHolder($scope.holder, function(data){
            $scope.change_holder = data;
            if ($scope.change_holder == "Error: no cow catch found"){
                $("#error_holder").show();
                $("#success_holder").hide();
            } else{             
                $("#success_holder").show();
                $("#error_holder").hide();
            }               
        });             
    }    

});

// Angular Factory
app.factory('appFactory', function($http){
    
    var factory = {};

    // For traceability system
    factory.queryAllCow = function(callback){

        $http.get('/get_all_cow/').success(function(output){
            callback(output)
        });
    }

    factory.queryCow = function(id, callback){
        $http.get('/get_cow/'+id).success(function(output){
            callback(output)
        });
    }

    factory.registerCow = function(data, callback){

        var cow = data.trace_id + "-" + data.cow_birthday + "-" + data.cow_category + "-" + data.cow_sex + "-" +
            data.owner + "-" + data.register_category + "-" + data.register_date + "-" + data.owner_address;

        $http.get('/register_cow/'+cow).success(function(output){
            callback(output)
        });
    }

    factory.updateSlaughterInfoCow = function(data, callback){

        var slaughter_info = data.trace_id + "-" + data.slaughter_house + "-" + data.slaughter_date + "-" +
            data.cow_result + "-" + data.cow_weight + "-" + data.cow_grade + "-" + data.slaughter_company;

        $http.get('/update_slaughter_info_cow/'+slaughter_info).success(function(output){
            callback(output)
        });
    }

    factory.updatePackageInfoCow = function(data, callback){

        var package_info = data.trace_id + "-" + data.company + "-" + data.company_address + "-" +
            data.cow_part + "-" + data.package_amount + "-" + data.package_date;

        $http.get('/update_package_info_cow/'+package_info).success(function(output){
            callback(output)
        });
    }

    factory.recordCow = function(data, callback){

        data.location = data.longitude + ", "+ data.latitude;

        var cow = data.id + "-" + data.location + "-" + data.timestamp + "-" + data.holder + "-" + data.vessel;

        $http.get('/add_cow/'+cow).success(function(output){
            callback(output)
        });
    }

    factory.changeHolder = function(data, callback){

        var holder = data.id + "-" + data.name;

        $http.get('/change_holder/'+holder).success(function(output){
            callback(output)
        });
    }

    return factory;
});