import React, {useEffect, useReducer, useState} from 'react';
import { createHttpClient } from '../../http/client';
import { Loader} from '../Loader';
import { Alert } from '../Alert';
import { EmailForm } from './EmailForm';
import { MPlanForm } from './MPlanForm';
import { NotificationTimeSpanForm } from './NotificationTimeSpanForm';

/* TODO
- (done)Define SettingPage as functional component
- (done)Use SettingPage in setting.js and build this SettingPage step by step
- (WIP)Manage states with useState or useReducer hook
  - For instance, alert state should be define with useReducer maybe
  - https://ja.reactjs.org/docs/hooks-reference.html#usereducer
  - https://qiita.com/makishy/items/bb014073d6e494b1b35f
  - https://qiita.com/ossan-engineer/items/c3853315f59dc20bc9dc
  - stateの管理はSettingPageで行い、useReducerでもらったdispatchをEmailFormとかに渡す。EMailFormではsubmitされたらAPIを呼び出し、最後にdispatchを呼んで親に伝える。EMAIL_SUBMITTEDを受け取ったらalertのstateを更新してAlertを出す
  - isLoading
  - alert (Add AlertState)
  - email
  - timeSpan
  - mPlan
*/

type State = { // TODO: Maybe rename to SettingPageState
  loading: boolean;
  email: string;
};

type Action = {
  type: ActionType;
  payload: State;
};

enum ActionType {
  FETCH_INIT = 1,
  FETCH_SUCCESS,
  FETCH_ERROR,
  ON_CHANGE_EMAIL,
  UPDATE_EMAIL,
}

const reducer: React.Reducer<State, Action> = (state: State, action: Action) => {
  switch (action.type) {
    case ActionType.FETCH_INIT:
      return {
        ...state,
        loading: true,
      };
    case ActionType.FETCH_SUCCESS:
      console.log(state);
      return {
        ...state,
        loading: false,
      };
    case ActionType.FETCH_ERROR:
      return {
        ...state,
        loading: false,
      };
    case ActionType.ON_CHANGE_EMAIL:
      return {
        ...state,
        email: action.payload.email,
      }
    case ActionType.UPDATE_EMAIL:
      return {
        ...state,
        email: action.payload.email,
      }
    default:
      throw Error('Invalid action');
  }
}

export const SettingPage: React.FC<{}> = () => {
  const initialState:State = {
    loading: false,
    email: '',
  };
  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    dispatch({
      type: ActionType.FETCH_INIT,
      payload: {...state},
    });
    const client = createHttpClient();
    client
      .get('/api/v1/me')
      .then((response) => {
        console.log(response.data);
        const timeSpans = response.data['notificationTimeSpans']
          ? response.data['notificationTimeSpans']
          : [];
        dispatch({
          type: ActionType.FETCH_SUCCESS,
          payload: {
            ...response.data,
          },
        });
        // this.setState({
        //   loading: false,
        //   userId: response.data['userId'],
        //   email: response.data['email'],
        //   timeSpan: {
        //     editable: this.state.timeSpan.editable,
        //     timeSpans: timeSpans,
        //   },
        //   mPlan: response.data['mPlan'],
        // });
      })
      .catch((error) => {
        console.log(error);
        dispatch({
          type: ActionType.FETCH_ERROR,
          payload: { ...state },
        })
        //this.handleShowAlert('danger', 'システムエラーが発生しました');
      });
  }, [])

  console.log('return JSX: state=', state);
  return (
    <div>
     <h1 className="page-title">設定</h1>
      {
        state.loading ?
          <Loader
            loading={state.loading}
            message={'Loading data ...'}
            css={'background: rgba(255, 255, 255, 0)'}
            size={50}
          /> :
          <>
            {/*<Alert*/}
            {/*  {...this.state.alert}*/}
            {/*  handleCloseAlert={this.handleHideAlert}*/}
            {/*/>*/}
            <EmailForm
              email={state.email}
              handleOnChange={this.handleOnChangeEmail}
              handleUpdateEmail={this.handleUpdateEmail}
            />
            <NotificationTimeSpanForm
              handleAdd={this.handleAddTimeSpan}
              handleDelete={this.handleDeleteTimeSpan}
              handleUpdate={this.handleUpdateTimeSpan}
              handleOnChange={this.handleOnChangeTimeSpan}
              handleSetEditable={this.handleSetTimeSpanEditable}
              {...this.state.timeSpan}
            />
            <MPlanForm {...this.state.mPlan} />
          </>
      }
    </div>
  );
}

// class SettingView extends MicroContainer {
//   constructor(props) {
//     super(props);
//     this.state = {
//       loading: false,
//       alert: {
//         visible: false,
//         kind: '',
//         message: '',
//       },
//       userId: 0,
//       email: '',
//       timeSpan: {
//         editable: false,
//         timeSpans: [], // {fromHour:23, fromMinute:0, toHour:23, toMinute:30}
//       },
//       mPlan: {
//         id: 0,
//         name: '',
//       },
//     };
//
//     this.handleShowAlert = this.handleShowAlert.bind(this);
//     this.handleHideAlert = this.handleHideAlert.bind(this);
//     this.handleOnChangeEmail = this.handleOnChangeEmail.bind(this);
//     this.handleUpdateEmail = this.handleUpdateEmail.bind(this);
//     this.handleAddTimeSpan = this.handleAddTimeSpan.bind(this);
//     this.handleDeleteTimeSpan = this.handleDeleteTimeSpan.bind(this);
//     this.handleUpdateTimeSpan = this.handleUpdateTimeSpan.bind(this);
//     this.handleOnChangeTimeSpan = this.handleOnChangeTimeSpan.bind(this);
//     this.handleSetTimeSpanEditable = this.handleSetTimeSpanEditable.bind(this);
//   }
//
//   componentDidMount() {
//     this.subscribe({
//       showAlert: this.handleShowAlert,
//       hideAlert: this.handleHideAlert,
//       // Email
//       onChangeEmail: this.handleOnChangeEmail,
//       updateEmail: this.handleUpdateEmail,
//       // TimeSpan
//       setTimeSpanEditable: this.handleSetTimeSpanEditable,
//       addTimeSpan: this.handleAddTimeSpan,
//       deleteTimeSpan: this.handleDeleteTimeSpan,
//       updateTimeSpan: this.handleUpdateTimeSpan,
//       onChangeTimeSpan: this.handleOnChangeTimeSpan,
//     });
//
//     this.setState({
//       loading: true,
//     });
//     this.fetch();
//   }
//
//   render() {
//     return (
//       <div>
//         <h1 className="page-title">設定</h1>
//         <Loadable
//           active={this.state.loading}
//           spinner={
//             <Loader type="Oval" color="#00BFFF" height="100" width="100" />
//           }
//           text="Loading data ..."
//           styles={{
//             overlay: (base) => ({
//               ...base,
//               background: 'rgba(255, 255, 255, 0)',
//               color: 'rgba(0, 0, 0, 1)',
//             }),
//           }}
//         >
//           <Alert
//             {...this.state.alert}
//             handleCloseAlert={this.handleHideAlert}
//           />
//           <EmailForm
//             email={this.state.email}
//             handleOnChange={this.handleOnChangeEmail}
//             handleUpdateEmail={this.handleUpdateEmail}
//           />
//           <NotificationTimeSpanForm
//             handleAdd={this.handleAddTimeSpan}
//             handleDelete={this.handleDeleteTimeSpan}
//             handleUpdate={this.handleUpdateTimeSpan}
//             handleOnChange={this.handleOnChangeTimeSpan}
//             handleSetEditable={this.handleSetTimeSpanEditable}
//             {...this.state.timeSpan}
//           />
//           <MPlanForm {...this.state.mPlan} />
//         </Loadable>
//       </div>
//     );
//   }
//
//   fetch() {
//     const client = createHttpClient();
//     client
//       .get('/api/v1/me')
//       .then((response) => {
//         console.log(response.data);
//         const timeSpans = response.data['notificationTimeSpans']
//           ? response.data['notificationTimeSpans']
//           : [];
//         this.setState({
//           loading: false,
//           userId: response.data['userId'],
//           email: response.data['email'],
//           timeSpan: {
//             editable: this.state.timeSpan.editable,
//             timeSpans: timeSpans,
//           },
//           mPlan: response.data['mPlan'],
//         });
//       })
//       .catch((error) => {
//         console.log(error);
//         this.setState({
//           loading: false,
//         });
//         this.handleShowAlert('danger', 'システムエラーが発生しました');
//       });
//   }
//
//   handleShowAlert(kind, message) {
//     this.setState({
//       alert: { visible: true, kind: kind, message: message },
//     });
//   }
//
//   handleHideAlert() {
//     this.setState({
//       alert: { visible: false },
//     });
//   }
//
//   handleOnChangeEmail(e) {
//     this.setState({ email: e.target.value });
//   }
//
//   handleUpdateEmail(email) {
//     const client = createHttpClient();
//     client
//       .post('/api/v1/me/email', {
//         email: email,
//       })
//       .then((response) => {
//         this.handleShowAlert('success', 'メールアドレスを更新しました！');
//       })
//       .catch((error) => {
//         console.log(error);
//         if (error.response.status === 400) {
//           this.handleShowAlert(
//             'danger',
//             '正しいメールアドレスを入力してください'
//           );
//         } else {
//           // TODO: external message
//           this.handleShowAlert('danger', 'システムエラーが発生しました');
//         }
//       });
//   }
//
//   handleSetTimeSpanEditable(value) {
//     this.setState({
//       timeSpan: {
//         editable: value,
//         timeSpans: this.state.timeSpan.timeSpans,
//       },
//     });
//   }
//
//   handleAddTimeSpan() {
//     if (this.state.timeSpan.timeSpans.length === 3) {
//       return;
//     }
//     this.setState({
//       timeSpan: {
//         editable: this.state.timeSpan.editable,
//         timeSpans: [
//           ...this.state.timeSpan.timeSpans,
//           { fromHour: 0, fromMinute: 0, toHour: 0, toMinute: 0 },
//         ],
//       },
//     });
//   }
//
//   handleDeleteTimeSpan(index) {
//     let timeSpans = this.state.timeSpan.timeSpans.slice();
//     if (index >= timeSpans.length) {
//       return;
//     }
//     timeSpans.splice(index, 1);
//     this.setState({
//       timeSpan: {
//         editable: this.state.timeSpan.editable,
//         timeSpans: timeSpans,
//       },
//     });
//   }
//
//   handleOnChangeTimeSpan(name, index, value) {
//     let timeSpans = this.state.timeSpan.timeSpans.slice();
//     timeSpans[index][name] = value;
//     this.setState({
//       timeSpan: {
//         editable: this.state.timeSpan.editable,
//         timeSpans: timeSpans,
//       },
//     });
//   }
//
//   handleUpdateTimeSpan() {
//     const timeSpans = [];
//     for (const timeSpan of this.state.timeSpan.timeSpans) {
//       for (const [k, v] of Object.entries(timeSpan)) {
//         timeSpan[k] = parseInt(v);
//       }
//       if (
//         timeSpan.fromHour === 0 &&
//         timeSpan.fromMinute === 0 &&
//         timeSpan.toHour === 0 &&
//         timeSpan.toMinute === 0
//       ) {
//         // Ignore zero value
//         continue;
//       }
//       timeSpans.push(timeSpan);
//     }
//
//     const client = createHttpClient();
//     client
//       .post('/api/v1/me/notificationTimeSpan', {
//         notificationTimeSpans: timeSpans,
//       })
//       .then((response) => {
//         this.handleShowAlert('success', 'レッスン希望時間帯を更新しました！');
//       })
//       .catch((error) => {
//         console.log(error);
//         if (error.response.status === 400) {
//           this.handleShowAlert(
//             'danger',
//             '正しいレッスン希望時間帯を選択してください'
//           );
//         } else {
//           // TODO: external message
//           this.handleShowAlert('danger', 'システムエラーが発生しました');
//         }
//       });
//
//     this.setState({
//       timeSpan: {
//         editable: false,
//         timeSpans: timeSpans,
//       },
//     });
//   }
// }
//
// ReactDOM.render(<SettingView />, document.getElementById('root'));
