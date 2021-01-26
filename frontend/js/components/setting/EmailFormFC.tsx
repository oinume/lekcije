import React from "react";

type Props = {
  email: string
  handleOnChange: () => void
  handleUpdateEmail: (email: string) => void
}

export const EmailFormFC: React.FC<Props> = ({email, handleOnChange, handleUpdateEmail}) => {
  return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-3">
            <label htmlFor="email" className="control-label">Email</label>
          </div>
          <div className="col-sm-7">
            <input
              type="email" className="form-control" name="email" id="email"
              placeholder="Email" required autoFocus autoComplete="on"
              value={email} onChange={handleOnChange}/>
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
            <button
              type="button"
              disabled={!email}
              className="btn btn-primary"
              onClick={() => handleUpdateEmail(email)}
            >
              変更
            </button>
          </div>
        </div>
      </form>
  );
}
