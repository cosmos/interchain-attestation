import clsx from 'clsx';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Link from "@docusaurus/core/lib/client/exports/Link";
import { useColorMode } from '@docusaurus/theme-common';
import Heading from "@theme/Heading";

import styles from './index.module.css';

function HomepageHeader() {
  const { colorMode } = useColorMode();
  const LogoSvg = colorMode === 'dark'
        ? require('@site/static/img/logo-with-name-and-tagline-dark.svg').default
        : require('@site/static/img/logo-with-name-and-tagline.svg').default;
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
          <LogoSvg className={styles.logoSvg} role="img"></LogoSvg>
      </div>
    </header>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={siteConfig.title}
      description="Interchain Attestation is a project to enable IBC everywhere.">
      <HomepageHeader />
      <main>
          <section className={styles.introSection}>
              <Heading as="h2">Enabling IBC everywhere</Heading>
              <p>
                  Interchain Attestation is a project to enable IBC everywhere.
                  <br />
                  It enables IBC when a light client is not available, by letting the receiving chain attest to the sender chain.
              </p>
              <div className="flex justify-center">
                  <div className={styles.buttons}>
                      <Link
                          className="button button--secondary button--lg"
                          to="/docs/intro">
                          Read more in the docs
                      </Link>
                  </div>
              </div>
          </section>
      </main>
    </Layout>
  );
}
