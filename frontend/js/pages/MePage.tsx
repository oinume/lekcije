import React, {useState} from 'react';
import {PageTitle} from '../components/PageTitle';
import {useGetMe} from '../hooks/useGetMe';
import {Loader} from '../components/Loader';
import {ErrorAlert} from '../components/ErrorAlert';
import {Teacher} from '../models/Teacher';

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
    <CreateForm/>
    <TeacherList/>
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

// https://getbootstrap.com/docs/4.4/components/spinners/
const CreateForm = () => (
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
);

const TeacherList = () => {
  const teachers = [
    new Teacher(28_660, 'hoge'),
    new Teacher(28_661, 'fuga'),
    new Teacher(28_662, 'aaaa'),
  ];

  const [checkedIds, setCheckedIds] = useState<number[]>([]);
  const handleCheckboxOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const targetId = Number.parseInt(event.target.value, 10);
    if (event.target.checked) {
      setCheckedIds([...checkedIds, targetId]);
    } else {
      setCheckedIds(checkedIds.filter(id => id !== targetId));
    }
  };

  return (
    <div id="followingTeachers">
      <form method="POST" action="/me/followingTeachers/delete">
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th scope="col" className="col-md-1"><button type="submit" className="btn btn-primary btn-sm">削除</button></th>
              <th scope="col" className="col-md-11">講師</th>
            </tr>
          </thead>
          <tbody>
            {teachers.map(t => <TeacherRow key={t.id} teacher={t} handleOnChange={handleCheckboxOnChange}/>)}
          </tbody>
        </table>
      </form>
    </div>
  );
};

type TeacherRowProps = {
  teacher: Teacher;
  handleOnChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
};

const TeacherRow = ({teacher, handleOnChange}: TeacherRowProps) => (
  <tr>
    <td className="col-md-1"><input type="checkbox" name="teacherIds" value={teacher.id} onChange={handleOnChange}/></td>
    <td className="col-md-8"><a href={`https://eikaiwa.dmm.com/teacher/index/${teacher.id}`} target="_blank" rel="noreferrer">{ teacher.name }</a></td>
  </tr>
);
