new Vue({
    el: '#app',
    data: {
        ws: null, // websocket
        newMsg: '', // holds new messages to be sent to the server
        chatContent: '', // list of chat content displayed
        email: null, // used for grabbing avatar
        username: null, // username
        joined: false // true if email + username not null
    },

    // handles initial setup
    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            // get avatar + parse emojis
            self.chatContent += '<div class="chip">' + '<img src="' + self.gravatarURL(msg.email) + '">' + msg.username + '</div>' + emojione.toImage(msg.message) + '</br'; 
            var element = document.getElementById('chat-messages');
            // auto scroll to bottom
            element.scrollTop = element.scrollHeight;
        });
    },

    // methods
    methods: {
        // send messages to server
        send: function() {
            // check message is not blank
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // sanitise by stripping html + js ?
                    }
                ));
                // reset newMsg
                this.newMsg = '';
            }
        },
        // ensure user enters email + username before sending messages
        join: function() {
            if (!this.email) {
                Materialize.toast('Please enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('Please enter a witty username', 2000);
                return
            }
            // sanitise inputs and set joined to true
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },
        // get gravatar URL
        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});