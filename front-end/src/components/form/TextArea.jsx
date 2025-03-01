const TextArea = (props) => {
  return (
    <div className="mb-3">
      <label htmlFor={props.name} className="form-label">
        {props.title}
      </label>
      <textarea
        className="form-control"
        id={props.name}
        value={props.value}
        name={props.name}
        onChange={props.onChange}
        rows={props.rows}
      />
      <div className={props.errorDiv}>{props.errorMsg}</div>
    </div>
  );
};

export default TextArea;
