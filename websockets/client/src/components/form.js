import React from "react";
import useForm from "../hooks/useForm";

const Form = ({ ws }) => {
  const { values, handleChange, handleSubmit } = useForm(send);
  function send() {
    ws.send(
      JSON.stringify({
        user: values.name,
        message: values.message
      })
    );
  }
  return (
    <form onSubmit={handleSubmit}>
      <div className="field">
        <label className="label">Name</label>
        <div className="control">
          <input
            className="input"
            type="text"
            name="name"
            onChange={handleChange}
            defaultValue={values.name}
          />
        </div>
      </div>
      <div className="field">
        <label className="label">Message</label>
        <div className="control">
          <input
            className="input"
            type="text"
            name="message"
            onChange={handleChange}
            defaultValue={values.message}
          />
        </div>
      </div>
      <button type="submit" className="button is-block is-info is-fullwidth">
        Send
      </button>
    </form>
  );
};

export default Form;
