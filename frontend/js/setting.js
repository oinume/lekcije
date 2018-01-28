import React from 'react';
import ReactDOM from 'react-dom';
import MicroContainer from 'react-micro-container';
import {createHttpClient} from './http';
import Alert from './components/Alert';
import Select from './components/Select';
import {sprintf} from 'sprintf-js';

class SettingView extends MicroContainer {

  constructor(props) {
    super(props);
    this.state = {
      alert: {
        visible: false,
        kind: '',
        message: '',
      },
      email: '',
      timeSpan: {
        editable: false,
        timeSpans: [], // {fromHour:23, fromMinute:0, toHour:23, toMinute:30}
      },
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

    this.fetch();
  }

  render() {
    return (
      <div>
        <Alert dispatch={this.dispatch} {...this.state.alert}/>
        <EmailForm dispatch={this.dispatch} value={this.state.email}/>
        <NotificationTimeSpanForm dispatch={this.dispatch} {...this.state.timeSpan}
        />
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
          email: response.data['email'],
          timeSpan: {
            editable: this.state.timeSpan.editable,
            timeSpans: timeSpans,
          },
        })
      })
      .catch((error) => {
        console.log(error);
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
        this.handleShowAlert('success', '通知対象の時間帯を更新しました！');
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 400) {
          this.handleShowAlert('danger', '正しい通知対象の時間帯を選択してください');
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

class EmailForm extends React.Component {

  constructor(props) {
    super(props);
    this.onChange = this.onChange.bind(this);
  }

  onChange(e) {
    this.props.dispatch('onChangeEmail', e.target.value);
  }

  render() {
    return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-2">
            <label htmlFor="email" className="control-label">Email address</label>
          </div>
          <div className="col-sm-8">
            <input
              type="email" className="form-control" name="email" id="email"
              placeholder="Email" required autoFocus autoComplete="on"
              value={this.props.value} onChange={this.onChange}/>
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
            <button
              type="button"
              disabled={!this.props.value}
              className="btn btn-primary"
              onClick={() => this.props.dispatch('updateEmail', this.props.value)}
            >
              変更
            </button>
          </div>
        </div>
      </form>
    );
  }
}

class NotificationTimeSpanForm extends React.Component {
  constructor(props) {
    super(props);
    this.onClickPlus = this.onClickPlus.bind(this);
    this.onChange = this.onChange.bind(this);
  }

  onClickPlus(e) {
    e.preventDefault();
    this.props.dispatch('addTimeSpan');
  }

  onClickMinus(e, index) {
    e.preventDefault();
    this.props.dispatch('deleteTimeSpan', index);
  }

  onChange(event) {
    const a = event.target.name.split('_');
    const name = a[0];
    const index = a[1];
    const timeSpans = this.props.timeSpans.slice();
    timeSpans[index][name] = event.target.value;
    this.props.dispatch('onChangeTimeSpan', name, index, event.target.value);
  }

  createTimeSpanRow(timeSpan, index) {
    let hourOptions = [];
    for (let i = 0; i <= 23; i++) {
      hourOptions.push({value: i, label: i})
    }
    let minuteOptions = [];
    for (const i of [0, 30]) {
      minuteOptions.push({value: i, label: sprintf('%02d', i)});
    }

    return (
      <p>
        <Select
          name={'fromHour_' + index}
          value={timeSpan.fromHour}
          onChange={this.onChange}
          options={hourOptions}
          className="custom-select mr-sm-2"
        />時
        &nbsp;
        <Select
          name={'fromMinute_' + index}
          value={timeSpan.fromMinute}
          onChange={this.onChange}
          options={minuteOptions}
          className="custom-select mr-sm-2"
        />分

        &nbsp; 〜 &nbsp;&nbsp;

        <Select
          name={'toHour_' + index}
          value={timeSpan.toHour}
          onChange={this.onChange}
          options={hourOptions}
          className="custom-select mr-sm-2"
        />時
        &nbsp;
        <Select
          name={'toMinute_' + index}
          value={timeSpan.toMinute}
          onChange={this.onChange}
          options={minuteOptions}
          className="custom-select mr-sm-2"
        />分

        <button type="button" className="btn btn-link btn-xs" onClick={() => this.onClickPlus(event)}>
          <span className="glyphicon glyphicon-plus-sign"/>
        </button>
        <button type="button" className="btn btn-link btn-xs" onClick={() => this.onClickMinus(event, index)}>
          <span className="glyphicon glyphicon-minus-sign"/>
        </button>
      </p>
    );
  }

  render() {
    let content = [];
    if (this.props.editable) {
      if (this.props.timeSpans.length > 0) {
        this.props.timeSpans.map((timeSpan, i) => {
          content.push(this.createTimeSpanRow(timeSpan, i));
        });
      } else {
        this.props.dispatch('addTimeSpan');
      }
    } else {
      if (this.props.timeSpans.length > 0) {
        for (let timeSpan of this.props.timeSpans) {
          content.push(
            <p>{timeSpan.fromHour}:{sprintf('%02d', timeSpan.fromMinute)} 〜 {timeSpan.toHour}:{sprintf('%02d', timeSpan.toMinute)}</p>);
        }
      } else {
        content.push(
          <p>データがありません。編集ボタンで追加できます。</p>
        );
      }
    }

    let updateButton;
    if (this.props.editable) {
      updateButton =
        <button
          type="button" className="btn btn-primary"
          onClick={() => this.props.dispatch('updateTimeSpan', false)}
        >
          更新
        </button>;
    } else {
      updateButton =
        <button
          type="button" className="btn btn-default"
          onClick={() => this.props.dispatch('setTimeSpanEditable', true)}
        >
          編集
        </button>;
    }

    return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-2">
            <label htmlFor="notificationTimeSpan" className="control-label">通知対象の時間帯</label>
            <a href="https://lekcije.amebaownd.com/posts/3614832" target="_blank">
              <span class="glyphicon glyphicon-question-sign" aria-hidden="true"></span>
            </a>
          </div>
          <div className="col-sm-8">
            {content}
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
            {updateButton}
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
