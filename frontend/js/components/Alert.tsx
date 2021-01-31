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
  handleCloseAlert: () => void
};

export const Alert: React.FC<Props> = ({kind, message, visible, handleCloseAlert}) => {
  if (visible) {
    return (
      <div className={"alert alert-" + kind} role="alert">
        <button
          type="button"
          className="close"
          onClick={() => handleCloseAlert()}
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
