import React from 'react';
import {PageTitle} from '../components/PageTitle';
import {useGetMe} from '../hooks/useGetMe';
import {Loader} from '../components/Loader';
import {ErrorAlert} from '../components/ErrorAlert';

export const MePage = () => {
  const getMeResult = useGetMe({});

  return (
    <div id="followingForm">
      <PageTitle>フォローしている講師</PageTitle>
      {
        getMeResult.isLoading || getMeResult.isIdle
          ? <Loader loading={getMeResult.isLoading} message="Loading data ..." css="background: rgba(255, 255, 255, 0)" size={50}/>
          : <MeContent showTutorial={getMeResult.data!.showTutorial}/>
      }
      {getMeResult.isError ? <ErrorAlert message={getMeResult.error.message}/> : <div/>}
    </div>
  );
};

type MeContentProps = {
  showTutorial: boolean;
};

// Help URL
// https://lekcije.amebaownd.com/posts/{{ if .IsUserAgentPC }}2044879{{ end }}{{ if .IsUserAgentSP }}1577091{{ end }}{{ if .IsUserAgentTablet }}1577091{{ end }}

const MeContent = ({showTutorial}: MeContentProps) => (
  <>
    {showTutorial ? <Tutorial/> : <div/>}
    <form method="POST" action="/me/followingTeachers/create">
      <p>
        講師のURLまたはIDを入力してフォローします<a href="https://lekcije.amebaownd.com/posts/2044879" rel="noreferrer" target="_blank"><i className="fas fa-question-circle button-help" aria-hidden="true"/></a><br/>
        <small><a href="https://eikaiwa.dmm.com/" rel="noreferrer" target="_blank">DMM英会話で講師を検索</a></small>
      </p>
      <div className="input-group mb-3">
        <input
          id="teacherIdsOrUrl"
          type="text"
          className="form-control"
          name="teacherIdsOrUrl"
          placeholder="https://eikaiwa.dmm.com/teacher/index/492/"
        />
        <button type="submit" className="btn btn-primary">送信</button>
      </div>
    </form>
  </>
);

const Tutorial = () => (
  <div className="alert alert-success alert-dismissible" role="alert">
    <button type="button" className="btn-close" data-bs-dismiss="alert" aria-label="Close"/>
    <h4><i className="bi bi-info-square-fill"/> 講師をフォローするには</h4>
    <ol>
      <li><a href="https://eikaiwa.dmm.com/list/" className="alert-link" target="_blank" rel="noreferrer">DMM英会話</a>でお気に入りの講師のページにアクセスしましょう</li>
      <li>講師のURLをコピーしましょう(<a href="https://lekcije.amebaownd.com/posts/1577091" className="alert-link" target="_blank" rel="noreferrer">ヘルプ</a>)</li>
      <li>URLを下の入力欄にペーストしてフォローしましょう</li>
      <li>フォローすると、その講師の空きレッスンがあった時にメールでお知らせします</li>
    </ol>
  </div>
);
