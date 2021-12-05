import React, {useState} from 'react';
import {toast} from 'react-toastify';
import {useMutation, useQueryClient} from 'react-query';
import {PageTitle} from '../components/PageTitle';
import {useGetMe} from '../hooks/useGetMe';
import {Loader} from '../components/Loader';
import {ErrorAlert} from '../components/ErrorAlert';
import {Teacher} from '../models/Teacher';
import {useListFollowingTeachers} from '../hooks/useListFollowingTeachers';
import {ToastContainer} from '../components/ToastContainer';
import {TwirpError, twirpRequest} from '../http/twirp';
import {queryKeyFollowingTeachers} from '../hooks/common';

export const MePage = () => {
  const getMeResult = useGetMe({});
  return (
    <div id="followingForm">
      <ToastContainer
        closeOnClick={false}
      />
      <PageTitle>フォローしている講師</PageTitle>
      {
        getMeResult.isLoading || getMeResult.isIdle
          ? <Loader isLoading={getMeResult.isLoading}/>
          : <MeContent showTutorial={getMeResult.data!.showTutorial}/>
      }
      {getMeResult.isError ? <ErrorAlert message={getMeResult.error.message}/> : <div/>}
    </div>
  );
};

type MeContentProps = {
  showTutorial: boolean; // eslint-disable-line react/boolean-prop-naming
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

const CreateForm = () => {
  const [teacherIdOrUrl, setTeacherIdOrUrl] = useState<string>('');
  const [submitDisabled, setSubmitDisabled] = useState<boolean>(true);

  const queryClient = useQueryClient();
  const createFollowingTeacherMutation = useMutation(
    async (teacherIdOrUrl: string): Promise<Response> => twirpRequest(
      '/twirp/api.v1.Me/CreateFollowingTeacher',
      JSON.stringify({teacherIdOrUrl}),
    ),
    {
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKeyFollowingTeachers);
        setTeacherIdOrUrl('');
        setSubmitDisabled(true);
        toast.success('講師をフォローしました！');
      },
      onError: (error: TwirpError) => {
        // eslint-disable-next-line @typescript-eslint/no-base-to-string, @typescript-eslint/restrict-template-expressions
        console.error(`createFollowingTeacherMutation.onError: err=${error}`);
        toast.error(`講師のフォローに失敗しました: ${error.message}`);
      },
    },
  );

  return (
    <form
      onSubmit={event => {
        event.preventDefault();
        createFollowingTeacherMutation.mutate(teacherIdOrUrl);
      }}
    >
      <p>
        講師のURLまたはIDを入力してフォローします<a href="https://lekcije.amebaownd.com/posts/2044879" rel="noreferrer" target="_blank"><i className="fas fa-question-circle button-help" aria-hidden="true"/></a><br/>
        <small><a href="https://eikaiwa.dmm.com/" rel="noreferrer" target="_blank">DMM英会話で講師を検索</a></small>
      </p>
      <div className="input-group mb-3">
        <input
          required
          autoFocus
          id="teacherIdsOrUrl"
          type="text"
          className="form-control"
          name="teacherIdsOrUrl"
          placeholder="https://eikaiwa.dmm.com/teacher/index/492/"
          value={teacherIdOrUrl}
          onChange={event => {
            event.preventDefault();
            setSubmitDisabled(event.currentTarget.value === '');
            setTeacherIdOrUrl(event.currentTarget.value);
          }}
        />
        <button
          type="submit"
          className="btn btn-primary"
          disabled={submitDisabled}
        >
          送信
        </button>
      </div>
    </form>
  );
};

const TeacherList = () => {
  const [checkedIds, setCheckedIds] = useState<number[]>([]);
  const [deleteSubmitDisabled, setDeleteSubmitDisabled] = useState<boolean>(true);

  const handleCheckboxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const targetId = Number.parseInt(event.target.value, 10);
    if (event.target.checked) {
      setCheckedIds([...checkedIds, targetId]);
      setDeleteSubmitDisabled(false);
    } else {
      const restIds = checkedIds.filter(id => id !== targetId);
      setCheckedIds(restIds);
      setDeleteSubmitDisabled(restIds.length === 0);
    }
  };

  const queryClient = useQueryClient();
  const deleteFollowingTeacherMutation = useMutation(
    async (teacherIds: number[]): Promise<Response> => twirpRequest(
      '/twirp/api.v1.Me/DeleteFollowingTeachers',
      JSON.stringify({teacherIds}),
    ),
    {
      onSuccess: async () => {
        await queryClient.invalidateQueries(queryKeyFollowingTeachers);
        toast.success('講師のフォローを解除しました');
        setDeleteSubmitDisabled(true);
      },
      onError: (error: TwirpError) => {
        // eslint-disable-next-line @typescript-eslint/no-base-to-string, @typescript-eslint/restrict-template-expressions
        console.error(`deleteFollowingTeacherMutation.onError: err=${error}`);
        toast.error(`講師のフォロー解除に失敗しました: ${error.message}`);
      },
    },
  );

  const result = useListFollowingTeachers({});
  if (result.isLoading || result.isIdle) {
    return <Loader isLoading={result.isLoading}/>;
  }

  if (result.isError) {
    return <ErrorAlert message={result.error.message}/>;
  }

  const {teachers} = result.data;

  return (
    <div id="followingTeachers">
      <form
        onSubmit={event =>{
          event.preventDefault();
          deleteFollowingTeacherMutation.mutate(checkedIds);
        }}>
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th scope="col" className="col-md-1">
                <button
                  type="submit"
                  className="btn btn-primary btn-sm"
                  disabled={deleteSubmitDisabled}
                >
                  削除
                </button>
              </th>
              <th scope="col" className="col-md-11">
                講師
              </th>
            </tr>
          </thead>
          <tbody>
            {teachers.map(t => <TeacherRow key={t.id} teacher={t} handleOnChange={handleCheckboxChange}/>)}
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
