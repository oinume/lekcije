import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';
import MicroContainer from 'react-micro-container';
import cookie from 'cookie';

class SettingView extends MicroContainer {

  constructor(props) {
    super(props);
    this.state = {
      email: ''
    };
  }

  componentDidMount() {
    this.subscribe({
      fetch: this.handleFetch,
      onChange: this.handleOnChange,
      update: this.handleUpdate,
    });

    this.handleFetch();
  }

  render() {
    return (
      <div>
        <EmailField dispatch={this.dispatch} value={this.state.email} />
      </div>
    );
  }

  handleFetch() {
    // TODO: move util
    const cookies = cookie.parse(document.cookie);
    const headers = {};
    if (cookies['apiToken']) {
      headers['Grpc-Metadata-Api-Token'] = cookies['apiToken'];
      headers['X-Api-Token'] = cookies['apiToken'];
    }
    axios.get('/api/v1/me/email', {
      'headers': headers
    })
      .then((response) => {
        this.setState({
          email: response.data['email'],
        })
      })
      .catch((error) => {
        console.log(error);
        alert('GET failed'); // TODO: show error
      });
  }

  handleOnChange(email) {
    this.setState({email: email})
  }

  handleUpdate(email) {
    alert('email is ' + email)
  }
}

class EmailField extends React.Component {

//{{ template "_flashMessage.html" . }}
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
          <input type="email" className="form-control" name="email" id="email" placeholder="Email" required autoFocus autoComplete="on" value={this.props.value} onChange={this.onChange} />
        </div>
        <button
          type="button"
          disabled={!this.props.value}
          className="btn btn-primary"
          onClick={() => this.props.dispatch('update', this.props.value)}>送信</button>
      </form>
    );
  }
}

ReactDOM.render(
  <SettingView />,
  document.getElementById('root')
);
