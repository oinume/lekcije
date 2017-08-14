import React from 'react';
import ReactDOM from 'react-dom';

class SettingView extends React.Component {

  constructor() {
    super();
    this.state = {
      email: 'oinume@gmail.com'
    };
  }

  componentDidMount() {
    // ここでGET /api/v1/@me/setting を呼んでデータを取得する
  }

  render() {
    return (
      <div>
        <EmailField value={this.state.email} />
      </div>
    );
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
          <input type="email" className="form-control" name="email" id="email" placeholder="Email"required autoFocus autoComplete="on" value={this.props.value} />
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
