var app = angular.module('OtpApp', ['ngMaterial']);

// Intercept POST requests, convert to standard form encoding
app.config(['$httpProvider', function ($httpProvider) {
	$httpProvider.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
	$httpProvider.defaults.transformRequest.unshift(function (data, headersGetter) {
		var key, result = [];
		
		for (key in data) {
			if (data.hasOwnProperty(key)) {
				result.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
			}
		}
		return (result.join('&'));
	});
}]);

app.controller('OtpCtrl', ['$scope', '$http', '$interval', '$timeout', '$mdDialog', function($scope, $http, $interval, $timeout, $mdDialog) {

	$scope.newOtp = { name : "", key : "", time : 30, digits : 6 };
	$scope.otps = [ ];

	$http.get('/otp').success(function(data) {
	    $scope.otps = data;
	    $interval(refresh, 1000);
		refresh()
	}).error(handleError);

	$scope.removeOtp = function (otp) {
		$http.delete('/otp/' + otp.Id).success(function(data, status, headers, config) {
			$scope.otps.splice($scope.otps.indexOf(otp), 1);
		}).error(handleError);
	};

	$scope.addOtp = function (newOtp) {
		$http.post('/otp', newOtp).success(function(data, status, headers, config) {
			$scope.otps.push(data)
			refresh()
		}).error(handleError);

		$scope.newOtp = { name : "", key : "", time : 30, digits : 6 };
	};

	var handleError = function (data, status) {
		if (status === 401) { // if we're not logged in
			return window.location.reload();
		}

		$mdDialog.show(
			$mdDialog.alert()
				.parent(angular.element(document.body))
				.title('Oops there is an error (' + status + ')')
				.content(data.msg)
				.ariaLabel('Error Dialog')
				.ok('OK')
		);
	}

	var updateOtp = function (index) {
		$http.get('/otp/' + $scope.otps[index].Id).success(function(data) {
			if (data.Totp === $scope.otps[index].Totp) {
				return ($timeout(function () {
					updateOtp(index);
				}, 100));
			}
			data.Rest = data.Time;
			data.Percent = 100;
		    $scope.otps[index] = data;
		}).error(handleError);
	}

	var refresh = function () {
		var epoch = Math.round(new Date().getTime() / 1000.0);
		var otps = $scope.otps;

		for (var i = 0, l = otps.length; i < l; ++i) {
			var timeStep = otps[i].Time;
            var rest = timeStep - (epoch % timeStep)
            var percent = Math.round(100 * rest / timeStep);

            if (epoch % timeStep === 0) {
            	updateOtp(i);
            } else {
            	otps[i].Rest = rest;
            	otps[i].Percent = percent;
            }
		}
	};

}]);