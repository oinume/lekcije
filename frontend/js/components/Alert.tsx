import React, {useState} from 'react';

type Props = {
  kind: string;
  message: string;
};

export const Alert: React.FC<Props> = ({kind, message}) => {
  const [visible] = useState<boolean>(true);
  if (visible) {
    return (
      <div className={'alert alert-dismissible alert-' + kind} role="alert">
        <button type="button" className="btn-close" data-bs-dismiss="alert" aria-label="Close"/>
        {message}
      </div>
    );
  }

  return <div/>;
};
