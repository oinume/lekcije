import React from "react";

type Props = {
  name: string
};

export const MPlanFormFC: React.FC<Props> = ({name}) => {
    return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-3">
            <label htmlFor="plan" className="control-label">プラン</label>
          </div>
          <div className="col-sm-7">
            <p>{name}</p>
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
          </div>
        </div>
      </form>
    );
}
