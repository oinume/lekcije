import {Alert} from "./Alert";
import React from "react";

type Props = {
  message: string;
  isInternal?: boolean;
};

export const ErrorAlert = (props: Props) => {
  let message = 'エラーが発生しました: ';
  if (props.isInternal) {
    message = 'システムエラーが発生しました: ';
  }
  return <Alert kind="danger" message={message + props.message}/>;
}
