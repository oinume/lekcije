import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';
import MicroContainer from 'react-micro-container';

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
      update: this.handleUpdate,
    });

    this.handleFetch();
  }

  render() {
    return (
      <div>
        <EmailField value={this.state.email} />
      </div>
    );
  }

  handleFetch() {
    axios.get('/api/v1/setting/email')
      .then((response) => {
        this.setState({
          email: response.data['email'],
        })
      })
      .catch((error) => {
        console.log(error);
        alert('GET failed'); // TODO: show error and disable update button
      });
  }

  handleUpdate() {
  }
}

class EmailField extends React.Component {

//{{ template "_flashMessage.html" . }}
  constructor() {
    super();
  }

  render() {
    return (
      <form method="POST" action="/me/setting/update">
        <div className="form-group">
          <label htmlFor="email">Email address</label>
          <input type="email" className="form-control" name="email" id="email" placeholder="Email" required autoFocus autoComplete="on" value={this.props.value} />
        </div>
        <button type="submit" className="btn btn-primary">送信</button>
      </form>
    );
  }
}

ReactDOM.render(
  <SettingView />,
  document.getElementById('root')
);
