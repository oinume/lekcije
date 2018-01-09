import React from 'react';
import ReactDOM from 'react-dom';
import MicroContainer from 'react-micro-container';
import {createClient} from './http';
import Alert from './components/Alert';
import Select from './components/Select';

class SettingView extends MicroContainer {

  constructor(props) {
    super(props);
    this.state = {
      email: '',
      alert: {
        visible: false,
        kind: '',
        message: '',
      },
    };
  }

  componentDidMount() {
    this.subscribe({
      fetch: this.handleFetch,
      onChange: this.handleOnChange,
      update: this.handleUpdate, // TODO: rename to updateEmail or Move to EmailFormContainer
      showAlert: this.handleShowAlert,
      hideAlert: this.handleHideAlert,
    });

    this.handleFetch();
  }

  render() {
    const a = this.state.alert;
    return (
      <div>
        <Alert dispatch={this.dispatch} visible={a.visible} kind={a.kind} message={a.message}/>
        <EmailForm dispatch={this.dispatch} value={this.state.email}/>
        <NotificationTimeSpanFormContainer rootDispatch={this.dispatch}/>
      </div>
    );
  }

  handleFetch() {
    // TODO: ここでNotificationTimeSpanも取得
    const client = createClient();
    client.get('/api/v1/me/email')
      .then((response) => {
        this.setState({
          email: response.data['email'],
        })
      })
      .catch((error) => {
        console.log(error);
        this.handleShowAlert('danger', 'システムエラーが発生しました');
      });
  }

  handleOnChange(email) {
    this.setState({email: email})
  }

  handleUpdate(email) {
    //alert('email is ' + email);
    const client = createClient();
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
}

class EmailForm extends React.Component {

  constructor(props) {
    super(props);
    this.onChange = this.onChange.bind(this);
  }

  onChange(e) {
    this.props.dispatch('onChange', e.target.value);
  }

  render() {
    return (
      <form className="form-horizontal">
        <div className="form-group">
          <label htmlFor="email" className="col-sm-2 control-label">Email address</label>
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
              onClick={() => this.props.dispatch('update', this.props.value)}
            >
              変更
            </button>
          </div>
        </div>
      </form>
    );
  }
}

class NotificationTimeSpanFormContainer extends MicroContainer {

  constructor(props) {
    super(props);
    this.state = {
      editable: false,
      timeSpans: [], // {fromHour:23, fromMinutes:0, toHour:23, toMinutes:30}
    };
  }

  componentDidMount() {
    this.subscribe({
      setEditable: this.handleSetEditable,
      add: this.handleAdd,
      update: this.handleUpdate,
      onChange: this.handleOnChange,
    });

    this.fetchTimeSpans();
  }

  fetchTimeSpans() {
    let timeSpans = [
      {fromHour:0, fromMinutes:0, toHour:0, toMinutes:0}
    ];
    this.setState({
      timeSpans: timeSpans,
    });
  }

  handleSetEditable(value) {
    this.setState({editable: value})
  }

  handleAdd() {
    this.setState({
      timeSpans: [...this.state.timeSpans, {fromHour:0, fromMinutes:0, toHour:0, toMinutes:0}]
    });
  }

  handleOnChange(name, index, value) {
    let timeSpans = this.state.timeSpans;
    timeSpans[index][name] = value;
    this.setState({
      timeSpans: timeSpans,
    });
  }

  handleUpdate() {
    // TODO: api call
    this.setState({
      editable: false,
    });
  }

  render() {
    return <NotificationTimeSpanForm dispatch={this.dispatch} {...this.state}/>;
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
    this.props.dispatch('add');
  }

  onChange(event) {
    const a = event.target.name.split('_');
    const name = a[0];
    const index = a[1];
    const timeSpans = this.props.timeSpans.slice();
    timeSpans[index][name] = event.target.value;
    this.props.dispatch('onChange', name, index, event.target.value);
  }

  createTimeSpanContent(timeSpan, index) {
    let hourOptions = [];
    for (let i = 0; i <= 23; i++) {
      hourOptions.push({value: i, label: i})
    }
    let minuteOptions = [];
    for (const i of [0, 30]) {
      minuteOptions.push({value:i, label: i});
    }

    return (
      <div className="col-sm-8">
        <Select
          name={'fromHour_' + index}
          value={timeSpan.fromHour}
          onChange={this.onChange}
          options={hourOptions}
          className="custom-select mr-sm-2"
        />時
        &nbsp;
        <Select
          name={'fromMinutes_' + index}
          value={timeSpan.fromMinutes}
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
          name={'toMinutes_' + index}
          value={timeSpan.toMinutes}
          onChange={this.onChange}
          options={minuteOptions}
          className="custom-select mr-sm-2"
        />分

        /* TODO: add row max to 3 */
        <a href="" onClick={() => this.onClickPlus(event)}>
          <span className="glyphicon glyphicon-plus-sign"/>
        </a>
      </div>
    );
  }

  render() {
    let content;
    if (this.props.editable) {
      content = [];
      this.props.timeSpans.map((timeSpan, i) => {
        content.push(this.createTimeSpanContent(timeSpan, i));
      });
    } else {
      content =
        <div className="col-sm-8">
          <p>12:40 〜 23:50</p>
          <p>12:40 〜 23:50</p>
          <p>12:40 〜 23:50</p>
        </div>;
    }

    let updateButton;
    if (this.props.editable) {
      updateButton =
        <button
          type="button" className="btn btn-primary"
          onClick={() => this.props.dispatch('update', false)}
        >
          更新
        </button>;
    } else {
      updateButton =
        <button
          type="button" className="btn btn-default"
          onClick={() => this.props.dispatch('setEditable', true)}
        >
          編集
        </button>;
    }

    return (
      <form className="form-horizontal">
        <div className="form-group">
          <label htmlFor="notificationTimeSpan" className="col-sm-2 control-label">通知対象の時間帯</label>
          {content}
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
