import React from "react";

// TODO: interface or type alias?
// interface Props {
//     kind: string
//     message: string
//     visible: boolean
// }

type Props = {
  kind: string
  message: string
  visible: boolean
  handleHideAlert: () => void
};

export const AlertFC: React.FC<Props> = ({kind, message, visible, handleHideAlert}) => {
  if (visible) {
    return (
      <div className={"alert alert-" + kind} role="alert">
        <button
          type="button"
          className="close"
          onClick={() => handleHideAlert()}
        >
          &times;
        </button>
        {message}
      </div>
    );
  } else {
    return (
      <div/>
    );
  }
}
