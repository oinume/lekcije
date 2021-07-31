import React, { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from 'react-query';
import { createHttpClient } from '../../http/client';
import { sendRequest } from '../../http/fetch';
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
  type UpdateMeEmailResult = {};
  const updateMeEmailMutation = useMutation(
    (email: string): Promise<UpdateMeEmailResult> => {
      return sendRequest(
        '/twirp/api.v1.User/UpdateMeEmail',
        JSON.stringify({
          // TODO: Use proto generated code
          email: email,
        })
      );
    },
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

  console.log('BEFORE useQuery');
  const { isLoading, isIdle, error, data } = useQuery<GetMeResult, Error>(
    queryKeyMe,
    async () => {
      console.log('BEFORE fetch');
      const response = await sendRequest('/twirp/api.v1.User/GetMe', '{}');
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
      console.log('----- data -----');
      console.log(data);
      return data as GetMeResult;
    },
    {
      retry: 0,
    }
  );
  console.log('AFTER useQuery: isLoading = %s', isLoading);

  if (isLoading || isIdle) {
    // TODO: Loaderコンポーネントの子供にフォームのコンポーネントをセットして、フォームは出すようにする
    return (
      <Loader loading={isLoading} message={'Loading data ...'} css={'background: rgba(255, 255, 255, 0)'} size={50} />
    );
  }

  if (error) {
    console.error('error = %s', error);
    return <Alert kind={'danger'} message={'システムエラーが発生しました。' + error.message} />;
  }

  const safeData = data as GetMeResult; // TODO: better name
  const email = emailState ?? safeData.email;
  const notificationTimeSpans = notificationTimeSpansState ?? safeData.notificationTimeSpans;

  const handleShowAlert = (kind: string, message: string) => {
    setAlert({ visible: true, kind: kind, message: message });
  };

  const handleHideAlert = () => {
    setAlert({ ...alert, visible: false });
  };

  const handleAddTimeSpan = () => {
    if (notificationTimeSpans.length === 3) {
      return;
    }
    setNotificationTimeSpansState([...notificationTimeSpans, { fromHour: 0, fromMinute: 0, toHour: 0, toMinute: 0 }]);
  };

  const handleDeleteTimeSpan = (index: number) => {
    let timeSpans = notificationTimeSpans.slice();
    if (index >= timeSpans.length) {
      return;
    }
    timeSpans.splice(index, 1);
    setNotificationTimeSpansState(timeSpans);
  };

  const handleOnChangeTimeSpan = (name: string, index: number, value: number) => {
    let timeSpans = notificationTimeSpans.slice();
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

    const client = createHttpClient();
    client
      .post('/twirp/api.v1.User/UpdateMeNotificationTimeSpan', {
        notificationTimeSpans: timeSpans,
      })
      .then((_) => {
        handleShowAlert('success', 'レッスン希望時間帯を更新しました！');
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 400) {
          handleShowAlert('danger', '正しいレッスン希望時間帯を選択してください');
        } else {
          // TODO: external message
          handleShowAlert('danger', 'システムエラーが発生しました');
        }
      });

    setNotificationTimeSpansState(timeSpans);
  };

  const showUpdateMeEmailAlert = () => {
    switch (updateMeEmailMutation.status) {
      case 'success':
        return <Alert kind={'success'} message={'メールアドレスを更新しました！'} />;
      case 'error':
        return <Alert kind={'danger'} message={updateMeEmailMutation.error as string} />;
      default:
        return <></>;
    }
  };

  return (
    <div>
      <h1 className="page-title">設定</h1>
      <>
        <ToggleAlert handleCloseAlert={handleHideAlert} {...alert} />
        {showUpdateMeEmailAlert()}
        <EmailForm
          email={email}
          handleOnChange={(e) => {
            setEmailState(e.currentTarget.value);
          }}
          handleUpdateEmail={(em): void => {
            updateMeEmailMutation.mutate(em);
          }}
        />
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
