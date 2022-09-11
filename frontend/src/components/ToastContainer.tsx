import React from 'react';
import type {
  ToastContainerProps} from 'react-toastify';
import {
  ToastContainer as ReactToastContainer,
  Zoom,
} from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css'; // eslint-disable-line import/no-unassigned-import

export const ToastContainer: React.FC<Partial<ToastContainerProps>> = props => (
  <ReactToastContainer
    pauseOnFocusLoss
    pauseOnHover
    hideProgressBar
    icon
    newestOnTop
    autoClose={3000}
    position="top-center"
    rtl={false}
    transition={Zoom}
    {...props}
  />
);

/*
Const ToastCloseButton = ({
  closeToast,
}: {
  closeToast: React.MouseEventHandler<HTMLButtonElement>;
}) => <button className="delete" onClick={closeToast}/>;

const getClassName = (type?: TypeOptions) => {
  switch (type) {
    case 'info':
      return 'is-info';
    case 'success':
      return 'is-success';
    case 'warning':
      return 'is-warning';
    case 'error':
      return 'is-danger';
    default:
      return 'is-info';
  }
};
*/
