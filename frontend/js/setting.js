import React from 'react';
import ReactDOM from 'react-dom';
import MicroContainer from 'react-micro-container';
import {createClient} from './http';
import Alert from './components/Alert';

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
      update: this.handleUpdate, // TODO: rename to updateEmail
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
      </div>
    );
  }

  handleFetch() {
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
      <form method="POST" action="/me/setting/update">
        <div className="form-group">
          <label htmlFor="email">Email address</label>
          <input
            type="email" className="form-control" name="email" id="email"
            placeholder="Email" required autoFocus autoComplete="on"
            value={this.props.value} onChange={this.onChange}/>
        </div>
        <button
          type="button"
          disabled={!this.props.value}
          className="btn btn-primary"
          onClick={() => this.props.dispatch('update', this.props.value)}
        >
          送信
        </button>
      </form>
    );
  }
}

ReactDOM.render(
  <SettingView/>,
  document.getElementById('root')
);
