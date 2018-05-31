import React from 'react';
import ReactDOM from 'react-dom';
import MicroContainer from 'react-micro-container';
import {createHttpClient} from './http';
import Alert from './components/Alert';
import EmailForm from './components/setting/EmailForm';
import NotificationTimeSpanForm from './components/setting/NotificationTimeSpanForm';
import Loadable from 'react-loading-overlay';

class SettingView extends MicroContainer {

  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      alert: {
        visible: false,
        kind: '',
        message: '',
      },
      userId: 0,
      email: '',
      timeSpan: {
        editable: false,
        timeSpans: [], // {fromHour:23, fromMinute:0, toHour:23, toMinute:30}
      },
      plan: {
        id: 1,
        name: 'フリー',
      }
    };
  }

  componentDidMount() {
    this.subscribe({
      showAlert: this.handleShowAlert,
      // Email
      hideAlert: this.handleHideAlert,
      onChangeEmail: this.handleOnChangeEmail,
      updateEmail: this.handleUpdateEmail,
      // TimeSpan
      setTimeSpanEditable: this.handleSetTimeSpanEditable,
      addTimeSpan: this.handleAddTimeSpan,
      deleteTimeSpan: this.handleDeleteTimeSpan,
      updateTimeSpan: this.handleUpdateTimeSpan,
      onChangeTimeSpan: this.handleOnChangeTimeSpan,
    });

    this.setState({
      loading: true,
    });
    this.fetch();
  }

  render() {
    return (
      <div>
        <h1 className="page-title">設定</h1>
        <Loadable
          active={this.state.loading}
          background='rgba(255, 255, 255, 0)'
          color='rgba(0, 0, 0, 1)'
          spinner={true}
          text='Loading data ...'
        >
          <Alert dispatch={this.dispatch} {...this.state.alert}/>
          <EmailForm dispatch={this.dispatch} value={this.state.email}/>
          <NotificationTimeSpanForm dispatch={this.dispatch} {...this.state.timeSpan}/>
          <PlanForm {...this.state.plan}/>
        </Loadable>
      </div>
    );
  }

  fetch() {
    const client = createHttpClient();
    client.get('/api/v1/me')
      .then((response) => {
        console.log(response.data);
        const timeSpans = response.data['notificationTimeSpans'] ? response.data['notificationTimeSpans'] : [];
        this.setState({
          loading: false,
          userId: response.data['userId'],
          email: response.data['email'],
          timeSpan: {
            editable: this.state.timeSpan.editable,
            timeSpans: timeSpans,
          },
        })
      })
      .catch((error) => {
        console.log(error);
        this.setState({
          loading: false,
        });
        this.handleShowAlert('danger', 'システムエラーが発生しました');
      });
  }

  handleShowAlert(kind, message) {
    this.setState({
      alert: {visible: true, kind: kind, message: message}
    })
  }

  handleHideAlert() {
    this.setState({
      alert: {visible: false}
    })
  }

  handleOnChangeEmail(email) {
    this.setState({email: email})
  }

  handleUpdateEmail(email) {
    //alert('email is ' + email);
    const client = createHttpClient();
    client.post('/api/v1/me/email', {
      email: email,
    })
      .then((response) => {
        this.handleShowAlert('success', 'メールアドレスを更新しました！');
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 400) {
          this.handleShowAlert('danger', '正しいメールアドレスを入力してください');
        } else {
          // TODO: external message
          this.handleShowAlert('danger', 'システムエラーが発生しました');
        }
      });
  }

  handleSetTimeSpanEditable(value) {
    this.setState({
      timeSpan: {
        editable: value,
        timeSpans: this.state.timeSpan.timeSpans,
      }
    })
  }

  handleAddTimeSpan() {
    if (this.state.timeSpan.timeSpans.length === 3) {
      return;
    }
    this.setState({
      timeSpan: {
        editable: this.state.timeSpan.editable,
        timeSpans: [...this.state.timeSpan.timeSpans, {fromHour: 0, fromMinute: 0, toHour: 0, toMinute: 0}],
      }
    });
  }

  handleDeleteTimeSpan(index) {
    let timeSpans = this.state.timeSpan.timeSpans.slice();
    if (index >= timeSpans.length) {
      return;
    }
    timeSpans.splice(index, 1);
    this.setState({
      timeSpan: {
        editable: this.state.timeSpan.editable,
        timeSpans: timeSpans,
      },
    });
  }

  handleOnChangeTimeSpan(name, index, value) {
    let timeSpans = this.state.timeSpan.timeSpans.slice();
    timeSpans[index][name] = value;
    this.setState({
      timeSpan: {
        editable: this.state.timeSpan.editable,
        timeSpans: timeSpans,
      }
    });
  }

  handleUpdateTimeSpan() {
    const timeSpans = [];
    for (const timeSpan of this.state.timeSpan.timeSpans) {
      for (const [k, v] of Object.entries(timeSpan)) {
        timeSpan[k] = parseInt(v);
      }
      if (timeSpan.fromHour === 0
        && timeSpan.fromMinute === 0
        && timeSpan.toHour === 0
        && timeSpan.toMinute === 0) {
        // Ignore zero value
        continue;
      }
      timeSpans.push(timeSpan);
    }

    const client = createHttpClient();
    client.post('/api/v1/me/notificationTimeSpan', {
      notificationTimeSpans: timeSpans,
    })
      .then((response) => {
        this.handleShowAlert('success', 'レッスン希望時間帯を更新しました！');
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 400) {
          this.handleShowAlert('danger', '正しいレッスン希望時間帯を選択してください');
        } else {
          // TODO: external message
          this.handleShowAlert('danger', 'システムエラーが発生しました');
        }
      });


    this.setState({
      timeSpan: {
        editable: false,
        timeSpans: timeSpans,
      }
    });
  }
}

class PlanForm extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-3">
            <label htmlFor="plan" className="control-label">プラン</label>
          </div>
          <div className="col-sm-7">
            <p>{this.props.name}</p>
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
            <button
              type="button"
              disabled={!this.props.value}
              className="btn btn-primary"
              onClick={() => {}}
            >
              変更
            </button>
          </div>
        </div>
      </form>
    );
  }

}

ReactDOM.render(
  <SettingView/>,
  document.getElementById('root')
);
