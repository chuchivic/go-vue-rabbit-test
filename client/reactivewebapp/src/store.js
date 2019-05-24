import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

var _client;
export default new Vuex.Store({
  strict: false,
  state: {
    machines : [
      {
        id:1,
        title: 'Flux condenser',
        status: 'operational'
      },
      {
        id:2,
        title: 'Space crafter',
        status: 'operational'
      },
      {
        id:3,
        title: 'Annoying Death star',
        status: 'operational'
      },
      {
        id:4,
        title: 'Super sonic sex machine',

        status: 'operational'
      },
      {
        id:5,
        title: 'Vendo Opel Corsa 2009 62 CV',
        status: 'operational'
      },
      {
        id:6,
        title: 'Alguien acerca al escenario una cerveza bien fria porfavor ? <3',
        status: 'operational'
      }
    ],
    socket: {
      isConnected: false,
      message: '',
      reconnectError: false,
    }
  },

  mutations: {
    SOCKET_ONOPEN (state, event)  {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true
      console.log('connected')
      _client = event.currentTarget;
    },
    SOCKET_ONCLOSE (state, event)  {
      state.socket.isConnected = false
    },
    SOCKET_ONERROR (state, event)  {
      console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message)  {
      console.log('Mensaje: ' + JSON.stringify(message))
      state.socket.message = message
      var machineId = message.id;
      var foundMachine = state.machines.find((elm) => elm.id == machineId);
      var sentState = message.newstate;
      foundMachine.status = sentState;
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
    SOCKET_statechange (state, payload) {
      var machineId = payload.id;
      var foundMachine = state.machines.find((elm) => elm.id == machineId);
      var sentState = payload.newstate;
      foundMachine.status = sentState;
    }
  },
  actions: {
    emitNewState(context, id) {
      let msg = {id: 5, newstate:"chorizo"}

      console.log(msg);
      Vue.prototype.$socket.sendObj(msg)
    }
  }
})
