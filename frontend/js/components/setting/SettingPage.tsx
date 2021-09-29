import React, { useState } from 'react';
import { useMutation, useQuery, useQueryClient, UseMutationResult } from 'react-query';
import { sendRequest, HttpError } from '../../http/fetch';
import { Loader } from '../Loader';
import { Alert } from '../Alert';
import { ToggleAlert } from '../ToggleAlert';
import { EmailForm } from './EmailForm';
import { NotificationTimeSpan, NotificationTimeSpanForm } from './NotificationTimeSpanForm';

const queryKeyMe = 'me';

type ToggleAlertState = {
  visible: boolean;
  kind: string;
  message: string;
};

type GetMeResult = {
  email: string;
  notificationTimeSpans: NotificationTimeSpan[];
};

type UpdateMeEmailResult = {};

type UpdateMeNotificationTimeSPanResult = {};

export const SettingPage: React.FC<{}> = () => {
  const [alert, setAlert] = useState<ToggleAlertState>({
    visible: false,
    kind: '',
    message: '',
  });
  const [emailState, setEmailState] = useState<string | undefined>(undefined);
  const [notificationTimeSpansState, setNotificationTimeSpansState] = useState<NotificationTimeSpan[] | undefined>(
    undefined
  );

  const queryClient = useQueryClient();
  // https://react-query.tanstack.com/guides/mutations
  const updateMeEmailMutation = useMutation(
    (email: string): Promise<UpdateMeEmailResult> =>
      sendRequest(
        '/twirp/api.v1.Me/UpdateEmail',
        JSON.stringify({
          // TODO: Use proto generated code
          email,
        })
      ),
    {
      onSuccess: () => {
        queryClient
          .invalidateQueries(queryKeyMe)
          .then((_) => {})
          .catch((e) => {
            console.error(e);
          });
      },
    }
  );

  const updateMeNotificationTimeSpanMutation = useMutation(
    (timeSpans: NotificationTimeSpan[]): Promise<UpdateMeNotificationTimeSPanResult> =>
      sendRequest(
        '/twirp/api.v1.Me/UpdateNotificationTimeSpan',
        JSON.stringify({
          notificationTimeSpans: timeSpans,
        })
      ),
    {
      onSuccess: () => {
        queryClient
          .invalidateQueries(queryKeyMe)
          .then((_) => {})
          .catch((e) => {
            console.error(e);
          });
      },
    }
  );

  // console.log('BEFORE useQuery<GetMeResult, Error>');
  const { isLoading, isIdle, error, data } = useQuery<GetMeResult, Error>(
    queryKeyMe,
    async () => {
      // console.log('BEFORE fetch');
      const response = await sendRequest('/twirp/api.v1.Me/GetMe', '{}');
      if (!response.ok) {
        // TODO: error
        type TwirpError = {
          code: string;
          msg: string;
        };
        const e: TwirpError = await response.json();
        throw new Error(`${response.status}:${e.msg}`);
      }
      const data = await response.json();
      // console.log('----- data -----');
      // console.log(data);
      return data as GetMeResult;
    },
    {
      retry: 0,
    }
  );
  // console.log('AFTER useQuery: isLoading = %s', isLoading);

  if (isLoading || isIdle) {
    // TODO: Loaderコンポーネントの子供にフォームのコンポーネントをセットして、フォームは出すようにする
    return <Loader loading={isLoading} message="Loading data ..." css="background: rgba(255, 255, 255, 0)" size={50} />;
  }

  if (error) {
    console.error('error = %s', error);
    return <Alert kind="danger" message={`システムエラーが発生しました。${error.message}`} />;
  }

  const safeData = data as GetMeResult; // TODO: better name
  const email = emailState ?? safeData.email;
  const notificationTimeSpans = notificationTimeSpansState ?? safeData.notificationTimeSpans;

  const handleHideAlert = () => {
    setAlert({ ...alert, visible: false });
  };

  const handleAddTimeSpan = () => {
    const maxTimeSpans = 3;
    if (notificationTimeSpans.length >= maxTimeSpans) {
      return;
    }
    setNotificationTimeSpansState([...notificationTimeSpans, { fromHour: 0, fromMinute: 0, toHour: 0, toMinute: 0 }]);
  };

  const handleDeleteTimeSpan = (index: number) => {
    const timeSpans = notificationTimeSpans.slice();
    if (index >= timeSpans.length) {
      return;
    }
    timeSpans.splice(index, 1);
    setNotificationTimeSpansState(timeSpans);
  };

  const handleOnChangeTimeSpan = (name: string, index: number, value: number) => {
    const timeSpans = notificationTimeSpans.slice();
    timeSpans[index][name as keyof NotificationTimeSpan] = value;
    setNotificationTimeSpansState(timeSpans);
  };

  const handleUpdateTimeSpan = () => {
    const timeSpans: NotificationTimeSpan[] = [];
    for (const timeSpan of notificationTimeSpans) {
      for (const [k, v] of Object.entries(timeSpan)) {
        timeSpan[k as keyof NotificationTimeSpan] = v;
      }
      if (timeSpan.fromHour === 0 && timeSpan.fromMinute === 0 && timeSpan.toHour === 0 && timeSpan.toMinute === 0) {
        // Ignore zero value
        continue;
      }
      timeSpans.push(timeSpan);
    }
    setNotificationTimeSpansState(timeSpans);
    // console.log('BEFORE updateMeNotificationTimeSpanMutation.mutate()');
    updateMeNotificationTimeSpanMutation.mutate(timeSpans);
  };

  return (
    <div>
      <h1 className="page-title">設定</h1>
      <>
        <ToggleAlert handleCloseAlert={handleHideAlert} {...alert} />
        <UseMutationResultAlert result={updateMeEmailMutation} name="メールアドレス" />
        <UseMutationResultAlert result={updateMeNotificationTimeSpanMutation} name="レッスン希望時間帯" />
        <EmailForm
          email={email}
          handleOnChange={(e) => {
            setEmailState(e.currentTarget.value);
          }}
          handleUpdateEmail={(em): void => {
            updateMeEmailMutation.mutate(em);
          }}
        />
        <div className="mb-3" />
        <NotificationTimeSpanForm
          handleAdd={handleAddTimeSpan}
          handleDelete={handleDeleteTimeSpan}
          handleUpdate={handleUpdateTimeSpan}
          handleOnChange={handleOnChangeTimeSpan}
          timeSpans={notificationTimeSpans}
        />
      </>
    </div>
  );
};

type UseMutationResultAlertProps = {
  result: UseMutationResult<any, any, any, any>;
  name: string;
};

// TODO: external component
const UseMutationResultAlert: React.FC<UseMutationResultAlertProps> = ({
  result,
  name,
}: UseMutationResultAlertProps) => {
  const e: HttpError = result.error as HttpError;
  // TODO: error handlingをまともにする
  switch (result.status) {
    case 'success':
      return <Alert kind="success" message={`${name}を更新しました！`} />;
    case 'error':
      return <Alert kind="danger" message={`${name}の更新に失敗しました。(${e.response.statusText})`} />;
    default:
      return <></>;
  }
};
