import React, { useState, useEffect } from "react";

const Form = ({ ws }) => {
  const [messages, setMessages] = useState([]);
  // const [listening, setListening] = useState(false);
  useEffect(() => {
    ws.addEventListener("message", e => {
      const msg = JSON.parse(e.data);
      setMessages([...messages, msg]);
    });
    //setListening(true);
  });

  return (
    <div>
      {messages.map(message => {
        return (
          <div>
            <b>{message.user}: </b>
            <span>{message.message}</span>
          </div>
        );
      })}
    </div>
  );
};

export default Form;
