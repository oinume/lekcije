import React from 'react';
import {PageTitle} from '../components/PageTitle';
import {useGetMe} from '../hooks/useGetMe';
import {Loader} from '../components/Loader';
import {ErrorAlert} from '../components/ErrorAlert';

export const MePage = () => {
  const {isLoading: isLoadingGetMe, isIdle: isIdleGetMe, error: errorGetMe} = useGetMe({});

  return (
    <div id="followingForm">
      <PageTitle>フォローしている講師</PageTitle>
      {isLoadingGetMe || isIdleGetMe ? <Loader loading={isLoadingGetMe} message="Loading data ..." css="background: rgba(255, 255, 255, 0)" size={50}/> : <MeContent/>}
      {errorGetMe ? <ErrorAlert message={errorGetMe.message}/> : <div/>}
    </div>
  );
};

const MeContent = () => (
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
