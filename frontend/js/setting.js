import React from 'react';
import ReactDOM from 'react-dom';
import MicroContainer from 'react-micro-container';
import {createClient} from './http';

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
        alert('GET failed'); // TODO: show error
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
        alert('POST success');
      })
      .catch((error) => {
        console.log(error);
        alert('POST failed'); // TODO: show error
      });
  }
}

class EmailForm extends React.Component {

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
          <input type="email" className="form-control" name="email" id="email" placeholder="Email" required autoFocus
                 autoComplete="on" value={this.props.value} onChange={this.onChange}/>
        </div>
        <button
          type="button"
          disabled={!this.props.value}
          className="btn btn-primary"
          onClick={() => this.props.dispatch('update', this.props.value)}>送信
        </button>
      </form>
    );
  }
}

ReactDOM.render(
  <SettingView/>,
  document.getElementById('root')
);
