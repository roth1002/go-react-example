/**
 * @jsx React.DOM
 */
'use strict';

var React       = global.React = require('react');
var ReactRouter = require('react-router-component');
var superagent  = global.superagent = require('superagent');

var Pages       = ReactRouter.Pages;
var Page        = ReactRouter.Page;
var NotFound    = ReactRouter.NotFound;
var Link        = ReactRouter.Link;

var MainPage = React.createClass({

  render: function() {
    return (
      <div className="MainPage">
        <h1>Hello, anonymous!</h1>
        <p><Link href="/users/doe">Login</Link></p>
      </div>
    );
  }
});

var UserPage = React.createClass({
  statics: {
    getUserInfo: function(username, cb) {
      superagent.get(
        '/api/v1/users/' + username,
        function(err, res) {
          cb(err, res ? res.body : null);
        });
    }
  },

  getInitialState: function() {
    var state, username = this.props.username;
    UserPage.getUserInfo(this.props.username, function(err, res){
      if (typeof state === 'undefined') {
        state = res
      } else {
        this.setState(res);
      };
    }.bind(this));
    state = state || {
      username: username,
      name: username.charAt(0).toUpperCase() + username.slice(1)
    };
    return state;
  },

  componentWillReceiveProps: function(nextProps) {
    if (this.props.username !== nextProps.username) {
      UserPage.getUserInfo(nextProps.username, function(err, info) {
        if (err) {
          throw err;
        }
        this.setState(info);
      }.bind(this));
    }
  },

  render: function() {
    var otherUser = this.props.username === 'doe' ? 'ivan' : 'doe';
    return (
      <div className="UserPage">
        <h1>Hello, {this.state.name}!</h1>
        <p>
          Go to <Link href={"/users/" + otherUser}>/users/{otherUser}</Link>
        </p>
        <p><Link href="/">Logout</Link></p>
      </div>
    );
  }
});

var NotFoundHandler = React.createClass({

  render: function() {
    return (
      <div>
        <p>Page not found</p>
        <p><Link href="/">Logout</Link></p>
      </div>
    );
  }
});

var App = React.createClass({

  render: function() {
    return (
      <html>
        <head>
          <link rel="stylesheet" href="/static/css/style.css" />
          <script src="/static/bundle.js" />
          <title>Go React Example</title>
        </head>
        <Pages className="App" path={this.props.path}>
          <Page path="/" handler={MainPage} />
          <Page path="/users/:username" handler={UserPage} />
          <NotFound handler={NotFoundHandler} />
        </Pages>
      </html>
    );
  }
});

module.exports = global.App = App;

if (typeof window !== 'undefined') {
  window.onload = function() {
    React.render(<App />, document);
  }
}
