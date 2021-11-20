import React, {useState} from 'react';
import {useMutation, useQueryClient, UseMutationResult} from 'react-query';
import {sendRequest, HttpError} from '../http/fetch';
import {Loader} from '../components/Loader';
import {Alert} from '../components/Alert';
import {ToggleAlert} from '../components/ToggleAlert';
import {EmailForm} from '../components/setting/EmailForm';
import {NotificationTimeSpan, NotificationTimeSpanForm} from '../components/setting/NotificationTimeSpanForm';
import {PageTitle} from '../components/PageTitle';
import {useGetMe} from '../hooks/useGetMe';

const queryKeyMe = 'me';

type ToggleAlertState = {
  visible: boolean;
  kind: string;
  message: string;
};

export const SettingPage: React.FC = () => {
  const [alert, setAlert] = useState<ToggleAlertState>({
    visible: false,
    kind: '',
    message: '',
  });
  const [emailState, setEmailState] = useState<string | undefined>(undefined);
  const [notificationTimeSpansState, setNotificationTimeSpansState] = useState<NotificationTimeSpan[] | undefined>(
    undefined,
  );

  const queryClient = useQueryClient();
  // https://react-query.tanstack.com/guides/mutations
  const updateMeEmailMutation = useMutation(
    async (email: string): Promise<Response> => sendRequest(
      '/twirp/api.v1.Me/UpdateEmail',
      JSON.stringify({
        // TODO: Use proto generated code
        email,
      }),
    ),
    {
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKeyMe);
      },
    },
  );

  const updateMeNotificationTimeSpanMutation = useMutation(
    async (timeSpans: NotificationTimeSpan[]): Promise<Response> => sendRequest(
      '/twirp/api.v1.Me/UpdateNotificationTimeSpan',
      JSON.stringify({
        notificationTimeSpans: timeSpans,
      }),
    ),
    {
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKeyMe);
      },
    },
  );

  // Console.log('BEFORE useGetMe');
  const {isLoading, isIdle, error, data} = useGetMe({});
  // Console.log('AFTER useGetMe: isLoading = %s', isLoading);

  if (isLoading || isIdle) {
    // TODO: Loaderコンポーネントの子供にフォームのコンポーネントをセットして、フォームは出すようにする
    return (
      <Loader loading={isLoading} message="Loading data ..." css="background: rgba(255, 255, 255, 0)" size={50}/>
    );
  }

  if (error) {
    console.error('useGetMe: error = %s', error);
    return <Alert kind="danger" message={'システムエラーが発生しました。' + error.message}/>;
  }

  const safeData = data!; // TODO: better name
  const email = emailState ?? safeData.email;
  const notificationTimeSpans = notificationTimeSpansState ?? safeData.notificationTimeSpans;

  const handleHideAlert = () => {
    setAlert({...alert, visible: false});
  };

  const handleAddTimeSpan = () => {
    const maxTimeSpans = 3;
    if (notificationTimeSpans.length >= maxTimeSpans) {
      return;
    }

    setNotificationTimeSpansState([...notificationTimeSpans, {fromHour: 0, fromMinute: 0, toHour: 0, toMinute: 0}]);
  };

  const handleDeleteTimeSpan = (index: number) => {
    const timeSpans = [...notificationTimeSpans];
    if (index >= timeSpans.length) {
      return;
    }

    timeSpans.splice(index, 1);
    setNotificationTimeSpansState(timeSpans);
  };

  const handleOnChangeTimeSpan = (name: string, index: number, value: number) => {
    const timeSpans = [...notificationTimeSpans];
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
    // Console.log('BEFORE updateMeNotificationTimeSpanMutation.mutate()');
    updateMeNotificationTimeSpanMutation.mutate(timeSpans);
  };

  return (
    <div>
      <PageTitle>設定</PageTitle>
      <ToggleAlert handleCloseAlert={handleHideAlert} {...alert}/>
      <UseMutationResultAlert result={updateMeEmailMutation} name="メールアドレス"/>
      <UseMutationResultAlert result={updateMeNotificationTimeSpanMutation} name="レッスン希望時間帯"/>
      <EmailForm
        email={email}
        handleOnChange={event => {
          setEmailState(event.currentTarget.value);
        }}
        handleUpdateEmail={(em): void => {
          updateMeEmailMutation.mutate(em);
        }}
      />
      <div className="mb-3"/>
      <NotificationTimeSpanForm
        handleAdd={handleAddTimeSpan}
        handleDelete={handleDeleteTimeSpan}
        handleUpdate={handleUpdateTimeSpan}
        handleOnChange={handleOnChangeTimeSpan}
        timeSpans={notificationTimeSpans}
      />
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
  const error: HttpError = result.error as HttpError;
  // TODO: error handlingをまともにする
  switch (result.status) {
    case 'success':
      return <Alert kind="success" message={name + 'を更新しました！'}/>;
    case 'error':
      return <Alert kind="danger" message={name + `の更新に失敗しました。(${error.response.statusText})`}/>;
    default:
      return <div/>;
  }
};
