import React from 'react';
import {
  ToastContainer as ReactToastContainer,
  ToastContainerProps,
  TypeOptions,
  Zoom,
} from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export const ToastContainer: React.FC<Partial<ToastContainerProps>> = (
  props
) => (
  <ReactToastContainer
    autoClose={3000}
    hideProgressBar={true}
    icon={true}
    newestOnTop={true}
    pauseOnFocusLoss
    pauseOnHover
    position="top-center"
    rtl={false}
    transition={Zoom}
    {...props}
  />
);

const ToastCloseButton = ({
                            closeToast,
                          }: {
  closeToast: React.MouseEventHandler<HTMLButtonElement>;
}) => <button className="delete" onClick={closeToast} />;

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

