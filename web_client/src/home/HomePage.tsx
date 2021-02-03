import { Container, makeStyles, Theme, Typography } from '@material-ui/core';

import { Page } from '../layouts';
import { NavBar } from '../core/bars';

const useStyles = makeStyles((theme: Theme) => ({
  root: {
    textAlign: 'center',
    marginTop: theme.spacing(2),
  },
  title: {
    marginTop: theme.spacing(1),
  },
  content: {
    marginTop: theme.spacing(1),
  },
}));

function HomePage() {
  const classes = useStyles();

  return (
    <Page title="Authx">
      <NavBar />
      <Container component="main" maxWidth="md" className={classes.root}>
        <Typography variant="h2" className={classes.title}>
          Authx
        </Typography>
        <Typography variant="body2" className={classes.content}>
          An auth service written in Go.
        </Typography>
      </Container>
    </Page>
  );
}

export default HomePage;
