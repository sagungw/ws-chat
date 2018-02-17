function render() {
    new Vue({
        el: '#app',
        data: {
            ws: null,
            privateWs: null,
            userListWs: null,

            username: username,
            messages: [],
            loggedInUsers: [],
            newMessage: '',

            privateChatSessions: []
        },
        created: function() {
            var self = this;
            this.ws = new WebSocket('ws://' + window.location.host + '/ws/group');
            this.ws.addEventListener('message', function(e) {
                var msg = JSON.parse(e.data);
                self.messages.push(msg)
                $('.chat-box').scrollTop($('.chat-box')[0].scrollHeight);
            });

            this.privateWs = new WebSocket('ws://' + window.location.host + '/ws/private');
            this.privateWs.addEventListener('message', function(e) {
                var message = JSON.parse(e.data);
                if (message.to != self.username && message.from != self.username) {
                    return;
                }

                var to = message.from;
                var from = message.to;

                var sess = null;
                if (to == self.username) {
                    for (var i = 0; i < self.privateChatSessions.length; i++) {
                        if(self.privateChatSessions[i].to == message.to) {
                            sess = self.privateChatSessions[i];
                            break;
                        }
                    }
                } else {
                    for (var i = 0; i < self.privateChatSessions.length; i++) {
                        if(self.privateChatSessions[i].to == message.from) {
                            sess = self.privateChatSessions[i];
                            break;
                        }
                    }
                }
                if(sess != null) {
                    sess.messages.push(message.message);
                } else {
                    self.privateChatSessions.push({
                        to: to,
                        from: from,
                        messages: [message.message]
                    });
                }
            });

            this.userListWs = new WebSocket('ws://' + window.location.host + '/ws/user');
            this.userListWs.addEventListener('message', function(e) {
                var message = JSON.parse(e.data);
                self.loggedInUsers = message;
            });

            this.userListWs.onopen = function() {
                self.userListWs.send(JSON.stringify({username: username}));
            };
        },
        methods: {
            sendGlobalMessage: function() {
                if (this.newMessage == undefined || this.newMessage == '') {
                    return;
                }

                this.ws.send(
                    JSON.stringify({
                        username: this.username,
                        content: this.newMessage
                    }
                ));
                this.newMessage = '';
            },
            startPrivate: function(to) {
                this.privateChatSessions.push({
                    to: to,
                    from: this.username,
                    messages: []
                });
            },
            sendMessage: function(to) {
                var sess = null;
                for (var i = 0; i < this.privateChatSessions.length; i++) {
                    if(this.privateChatSessions[i].to == to) {
                        sess = this.privateChatSessions[i];
                        break;
                    }
                }

                this.privateWs.send(
                    JSON.stringify({
                        to: to,
                        from: this.username,
                        message: {
                            username: this.username,
                            content: sess.newMessage
                        }
                    }
                ));

                sess.newMessage = '';
            }
        }
    });
}

var username = localStorage['username'];
if (username == undefined || username == '') {
    var chance = new Chance();
    localStorage['username'] = chance.name();
    username = localStorage['username']
}

render();
