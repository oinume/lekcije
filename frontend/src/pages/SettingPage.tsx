import React, {useState} from 'react';
import {useQueryClient} from '@tanstack/react-query';
import {toast} from 'react-toastify';
import {Loader} from '../components/Loader';
import {ErrorAlert} from '../components/ErrorAlert';
import {ToggleAlert} from '../components/ToggleAlert';
import {EmailForm} from '../components/setting/EmailForm';
import {NotificationTimeSpanForm} from '../components/setting/NotificationTimeSpanForm';
import {PageTitle} from '../components/PageTitle';
import {NotificationTimeSpanModel} from '../models/NotificatonTimeSpan';
import type {
  GetViewerWithNotificationTimeSpansQuery, NotificationTimeSpan, NotificationTimeSpanInput,
} from '../graphql/generated';
import {
  useGetViewerWithNotificationTimeSpansQuery, useUpdateNotificationTimeSpansMutation, useUpdateViewerMutation,
} from '../graphql/generated';
import type {GraphQLError} from '../http/graphql';
import {createGraphQLClient, toMessage} from '../http/graphql';
import {ToastContainer} from '../components/ToastContainer';

type ToggleAlertState = {
  isVisible: boolean;
  kind: string;
  message: string;
};

export const SettingPage: React.FC = () => {
  const [alert, setAlert] = useState<ToggleAlertState>({
    isVisible: false,
    kind: '',
    message: '',
  });
  const [emailState, setEmailState] = useState<string | undefined>(undefined);
  const [notificationTimeSpansState, setNotificationTimeSpansState] = useState<NotificationTimeSpanModel[] | undefined>(
    undefined,
  );

  // https://tanstack.com/query/v4/docs/guides/mutations
  const queryClient = useQueryClient();
  const graphqlClient = createGraphQLClient();
  const updateViewerMutation = useUpdateViewerMutation<GraphQLError>(graphqlClient, {
    async onSuccess() {
      await queryClient.invalidateQueries(useGetViewerWithNotificationTimeSpansQuery.getKey());
      toast.success('メールアドレスを更新しました！');
    },
    async onError(error) {
      console.error(error.response);
      toast.error(toMessage(error));
    },
  });

  const updateNotificationTimeSpansMutation = useUpdateNotificationTimeSpansMutation<GraphQLError>(graphqlClient, {
    async onSuccess() {
      await queryClient.invalidateQueries(useGetViewerWithNotificationTimeSpansQuery.getKey());
      toast.success('レッスン希望時間帯を更新しました！');
    },
    async onError(error) {
      console.error(error.response);
      toast.error(toMessage(error));
    },
  });

  const queryResult = useGetViewerWithNotificationTimeSpansQuery<GetViewerWithNotificationTimeSpansQuery, GraphQLError>(graphqlClient);
  if (queryResult.isLoading) {
    // TODO: Loaderコンポーネントの子供にフォームのコンポーネントをセットして、フォームは出すようにする
    return (
      <Loader isLoading={queryResult.isLoading}/>
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

  if (queryResult.error) {
    // eslint-disable-next-line @typescript-eslint/no-base-to-string, @typescript-eslint/restrict-template-expressions
    console.error(`getViewerQuery: error = ${queryResult.error}, type=${typeof queryResult.error}`);
    return <ErrorAlert message={toMessage(queryResult.error)}/>;
  }

  const email = emailState ?? queryResult.data.viewer.email;
  const notificationTimeSpans = notificationTimeSpansState ?? toModels(queryResult.data.viewer.notificationTimeSpans);

  const handleHideAlert = () => {
    setAlert({...alert, isVisible: false});
  };

  const handleAddTimeSpan = () => {
    const maxTimeSpans = 3;
    if (notificationTimeSpans.length >= maxTimeSpans) {
      return;
    }

    setNotificationTimeSpansState([...notificationTimeSpans, new NotificationTimeSpanModel(0, 0, 0, 0)]);
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
    const timeSpans: NotificationTimeSpanModel[] = [];
    for (const timeSpan of notificationTimeSpans) {
      for (const [k, v] of Object.entries<NotificationTimeSpan>(timeSpan)) {
        timeSpan[k as keyof NotificationTimeSpan] = Number(v);
      }

      if (NotificationTimeSpanModel.fromObject(timeSpan).isZero()) { // `timeSpan` is object somehow...
        // Ignore zero value
        continue;
      }

      timeSpans.push(timeSpan);
    }

    setNotificationTimeSpansState(timeSpans);
    // Console.log('BEFORE updateMeNotificationTimeSpanMutation.mutate()');
    updateNotificationTimeSpansMutation.mutate({
      input: {timeSpans},
    });
  };

  return (
    <div>
      <ToastContainer closeOnClick={false}/>
      <PageTitle>設定</PageTitle>
      <ToggleAlert handleCloseAlert={handleHideAlert} {...alert}/>
      <EmailForm
        email={email}
        handleOnChange={event => {
          setEmailState(event.currentTarget.value);
        }}
        handleUpdateEmail={(em): void => {
          updateViewerMutation.mutate({input: {email: em}});
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

const toModels = (timeSpans: NotificationTimeSpan[]): NotificationTimeSpanModel[] => timeSpans.map<NotificationTimeSpanModel>(o => NotificationTimeSpanModel.fromObject(o));
