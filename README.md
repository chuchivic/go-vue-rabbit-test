# Vuejs <3 Rabbitmq


Bassed on: https://github.com/Djangoum/Vue-3-rabbitmq


But the backend done in Go iestead Nodejs and the connection between back and frontend with vue-native-websocket instead socketio



Central and main piece of the system is the rabbitmq.

consumer.go is the main entry point of the application back end.

Also you can find producer.go, this script will send in a infinite loop random states, it's used to simulate a third application sending data to our state manager application.

then in client/reactivewebapp, you can find the front end code. 

If you launch both applications, and different backends you can test what was shown in the talk.

The main goal was to show how a web app with web sockets could be scaled using rabbitmq and how can we take advantage with vuejs/vuex reactivity features to give a proper structure to our front end code.

Hope you liked the talk <3

Please do not doubt to contact me if you have any problem or question. arielamorgarcia@gmail.com

Best Regards! 
