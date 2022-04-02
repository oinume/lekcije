import React, {useState} from 'react';
import {useMutation, useQueryClient} from 'react-query';
import {Loader} from '../components/Loader';
import {ErrorAlert} from '../components/ErrorAlert';
import {ToggleAlert} from '../components/ToggleAlert';
import {EmailForm} from '../components/setting/EmailForm';
import {NotificationTimeSpanForm} from '../components/setting/NotificationTimeSpanForm';
import {PageTitle} from '../components/PageTitle';
import {useGetMe} from '../hooks/useGetMe';
import {NotificationTimeSpan} from '../models/NotificatonTimeSpan';
import {queryKeyMe} from '../hooks/common';
import {TwirpError, twirpRequest} from '../http/twirp';
import {UseMutationResultAlert} from '../components/UseMutationResultAlert';
import {GetViewerQuery, useGetViewerQuery} from '../graphql/generated';
import {createGraphQLClient, GraphQLError} from '../http/graphql';

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
    async (email: string): Promise<Response> => twirpRequest(
      '/twirp/api.v1.Me/UpdateEmail',
      JSON.stringify({
        email,
      }),
    ),
    {
      async onSuccess() {
        await queryClient.invalidateQueries(queryKeyMe);
      },
    },
  );

  const updateMeNotificationTimeSpanMutation = useMutation(
    async (timeSpans: NotificationTimeSpan[]): Promise<Response> => twirpRequest(
      '/twirp/api.v1.Me/UpdateNotificationTimeSpan',
      JSON.stringify({
        notificationTimeSpans: timeSpans,
      }),
    ),
    {
      async onSuccess() {
        await queryClient.invalidateQueries(queryKeyMe);
      },
    },
  );

  const client = createGraphQLClient();
  const getViewerQuery = useGetViewerQuery<GetViewerQuery, GraphQLError>(client);
  if (getViewerQuery.isLoading || getViewerQuery.isIdle) {
    // TODO: Loaderコンポーネントの子供にフォームのコンポーネントをセットして、フォームは出すようにする
    return (
      <Loader isLoading={getViewerQuery.isLoading}/>
    );
  }
  // Console.log('BEFORE useGetMe');
  // const getMeResult = useGetMe({});
  // Console.log('AFTER useGetMe: isLoading = %s', isLoading);

  // if (getMeResult.isLoading || getMeResult.isIdle) {
  //   // TODO: Loaderコンポーネントの子供にフォームのコンポーネントをセットして、フォームは出すようにする
  //   return (
  //     <Loader isLoading={getMeResult.isLoading}/>
  //   );
  // }

  if (getViewerQuery.error) {
    // eslint-disable-next-line @typescript-eslint/no-base-to-string, @typescript-eslint/restrict-template-expressions
    console.error(`getViewerQuery: error = ${getViewerQuery.error}, type=${typeof getViewerQuery.error}`);
    return <ErrorAlert message={getViewerQuery.error.message}/>;
  }

  const email = emailState ?? getViewerQuery.data.viewer.email;
  const notificationTimeSpans = notificationTimeSpansState ?? getViewerQuery.data.viewer..notificationTimeSpans;

  const handleHideAlert = () => {
    setAlert({...alert, visible: false});
  };

  const handleAddTimeSpan = () => {
    const maxTimeSpans = 3;
    if (notificationTimeSpans.length >= maxTimeSpans) {
      return;
    }

    setNotificationTimeSpansState([...notificationTimeSpans, new NotificationTimeSpan(0, 0, 0, 0)]);
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
      for (const [k, v] of Object.entries<NotificationTimeSpan>(timeSpan)) {
        timeSpan[k as keyof NotificationTimeSpan] = Number(v);
      }

      if (NotificationTimeSpan.fromObject(timeSpan).isZero()) { // `timeSpan` is object somehow...
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
