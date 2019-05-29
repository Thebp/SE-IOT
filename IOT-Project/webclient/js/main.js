function RoomController($scope, $http, $timeout) {
    var url = "http://mndkk.dk:50002";

    var updateRooms = function(source=null) {
        $http.get(url + "/rooms").then(function(response) {
            rooms = response.data;
            rooms.sort((a,b) => (a.id > b.id) ? 1 : -1)
            $scope.rooms = rooms;
            if(source != null) {
                $scope.room = getRoom(source)
                $scope.colorValue = $scope.room.led_config.color
            }
        })
    }

    var getRoom = function(id) {
        for(room of $scope.rooms) {
            if(room.id == id) {
                return room;
            }
        }
        return null;
    }

    $scope.saveRoom = function() {
        if($scope.room.id == null) {
            $http.post(url + "/rooms", JSON.stringify($scope.room)).then(function(response) {
                updateRooms();
                $scope.room = {}
            });
        } else {
            $http.put(url + "/rooms/" + $scope.room.id, JSON.stringify($scope.room)).then(function(response) {
                updateRooms();
            })
        }
    }

    $scope.roomOnClick = function(id) {
        if($scope.moving_board != null) {
            $http.put(url + "/boards/" + $scope.moving_board, JSON.stringify({id:$scope.moving_board, room_id:id})).then(function(response) {
                updateRooms($scope.room.id);
            })
            $scope.moving_board = null;
        } else {
            $scope.room = getRoom(id);
        }
    }

    $scope.moveBoard = function(id) {
        $scope.moving_board = id;
        alert("Click on a room to move the board to it.")
    }

    $scope.locate = function(id) {
        $http.post(url + "/boards/" + id + "/ping", "")
    }

    $scope.newRoom = function() {
        $scope.room = {};
    }

    updateRooms();
}