import React, { Component } from "react";
import Form from "./components/form";
import Chat from "./components/chat";

import "./App.css";

class App extends Component {
  constructor(props) {
    super(props);
    this.ws = new WebSocket("ws://" + window.location.host + "/ws");
  }
  render() {
    return (
      <div className="App">
        <Chat ws={this.ws} />
        <Form ws={this.ws} />
      </div>
    );
  }
}

export default App;
